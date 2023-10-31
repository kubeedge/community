# Copyright 2021 The KubeEdge Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

from typing import List
from collections import deque

from promptulate.llms import ChatOpenAI, BaseLLM
from promptulate.agents import BaseAgent
from promptulate.output_formatter import OutputFormatter
from promptulate.schema import MessageSet, SystemMessage, UserMessage, BaseMessage
from promptulate.utils.color_print import print_text

from bot_description.prompt import (
    SYSTEM_PROMPT_TEMPLATE,
    GENERATE_PLAN_SYSTEM_PROMPT,
    ENVIRONMENT_SUMMARY_SYSTEM_PROMPT,
    SECURITY_CHECK_SYSTEM_PROMPT,
)
from bot_description.operator import Operator
from bot_description.sensor import Sensor
from bot_description.schema import (
    Task,
    GeneratePlanResponse,
    CommandResponse,
)


class RobotController:
    def __init__(self, operators: List[Operator]) -> None:
        self.operators: List[Operator] = operators

        if self.get_operator("stop") is None:
            raise ValueError("Please provide stop operator")

    def execute_action(self, action_name: str, *args, **kwargs) -> None:
        """Execute action by action_name."""
        operator = self.get_operator(action_name)
        print_text(f"[command] run {action_name} {str(args)} {str(kwargs)}", "green")
        operator.execute_action(*args, **kwargs)

    def get_operator(self, action_name: str) -> Operator:
        return next(
            (operator for operator in self.operators if operator.name == action_name),
            None,
        )


class RobotObserver:
    def __init__(self, sensors: List[Sensor]) -> None:
        self.sensors: List[Sensor] = sensors

    def get_data(self, sensor_name: str) -> None:
        sensor = self.get_sensor(sensor_name)
        return sensor.get_data()

    def get_sensor(self, sensor_name: str) -> Sensor:
        return next(
            (sensor for sensor in self.sensors if sensor.name == sensor_name),
            None,
        )

    def show_all_data(self):
        return "\n".join([sensor.get_data() for sensor in self.sensors])


class RobotAgent(BaseAgent):
    def __init__(
        self, controller: RobotController, observer: RobotObserver, *args, **kwargs
    ):
        super().__init__(*args, **kwargs)
        self.llm = ChatOpenAI(temperature=0.0)
        self.robot_controller: RobotController = controller
        self.robot_observer: RobotObserver = observer
        self.user_demand: str = ""
        self.pending_tasks: deque = deque()
        self.completed_tasks: deque = deque()

    def get_llm(self) -> BaseLLM:
        return self.llm

    def _run(self, prompt: str, *args, **kwargs) -> str:
        """Main loop for RobotAgent"""
        self.user_demand = prompt

        # generate plan
        self.generate_tasks()

        # build system prompt
        environment = self.llm(
            ENVIRONMENT_SUMMARY_SYSTEM_PROMPT.format(
                environment=self.robot_observer.show_all_data()
            )
        )
        system_prompt: str = SYSTEM_PROMPT_TEMPLATE.format(
            skills=str(self.robot_controller.operators),
            environment=environment,
            user_demand=self.user_demand,
            pending_tasks=f"{str(self.pending_tasks)}",
        )
        formatter = OutputFormatter(CommandResponse)
        instruction = formatter.get_formatted_instructions()
        system_prompt = f"{system_prompt}\n{instruction}"

        messages: MessageSet = MessageSet(
            messages=[
                SystemMessage(content=system_prompt),
            ]
        )

        counter = 0
        while True:
            self.robot_controller.execute_action("stop")

            llm_output: BaseMessage = self.llm.predict(messages)
            messages.add_message(llm_output)
            messages.add_user_message("next")

            try:
                # get current task
                # TODO add retry if failed
                task: Task = formatter.formatting_result(llm_output.content).task
                print_text(f"[task] current task: {task.dict()}", "blue")

                if task.name == "stop" or not self.pending_tasks:
                    self.robot_controller.execute_action("stop")
                    return

                # run task
                task = self.pending_tasks.popleft()
                self.completed_tasks.append(task)
                self.robot_controller.execute_action(
                    action_name=task.name, **task.parameters
                )

            except Exception as e:
                print_text(f"[error] {e}", "red")
                self.robot_controller.execute_action("stop")
                return

            counter += 1

    def generate_tasks(self):
        """Generate tasks based on user demand."""
        # build conversation
        formatter = OutputFormatter(GeneratePlanResponse)
        instruction: str = formatter.get_formatted_instructions()
        messages: MessageSet = MessageSet(
            messages=[
                SystemMessage(content=f"{GENERATE_PLAN_SYSTEM_PROMPT}\n{instruction}"),
                UserMessage(content=f"##User input:\n{self.user_demand}"),
            ]
        )

        llm_output: str = self.llm.predict(messages).content

        response: GeneratePlanResponse = formatter.formatting_result(llm_output)
        self.pending_tasks = deque(response.tasks)
        print_text(f"Generate plans: {self.pending_tasks}", "green")


def example():
    operators = []
    sensors = []
    controller = RobotController(operators)
    observer = RobotObserver(sensors)
    RobotAgent(controller=controller, observer=observer)
