from typing import Dict, Any, List, Callable, Optional
from pydantic import BaseModel, Field


class Operator:
    """Operator provide the ability to control the specified behavior.
    Developer should pass callback to implement the ability to control
    specified behavior."""

    def __init__(self, name: str, description: str, callback: Callable):
        self.name = name
        self.description = description
        self.callback = callback

    def execute_action(self, *args, **kwargs):
        self.callback(*args, **kwargs)


class Sensor:
    """Sensor provide the ability to receive specified information of robot,
    such as odom data. Passing callback to get data. The best implementation
    is convert original data to semantic information in callback function."""

    def __init__(self, name: str, description: str, callback: Callable):
        self.name = name
        self.description = description
        self.callback = callback

    def get_data(self, *args, **kwargs) -> str:
        self.callback(*args, **kwargs)


class Task(BaseModel):
    """RobotAgent can control the robot by executing task."""

    id: int = Field(description="Autoincrement task id")
    name: str = Field(description="task name")
    parameters: Optional[Dict[str, Any]] = Field(description="task parameters")
    reason: str = Field(description="Reason for task execution")


class GeneratePlanResponse(BaseModel):
    """Output response of LLM generating task sequences."""
    tasks: List[Task] = Field(description="task sequences")


class CommandResponse(BaseModel):
    """Output response of LLM control command."""
    task: Task = Field(description="control task")


def get_task_examples_for_generating_plan():
    """Get task examples for output format few-shot."""
    return [
        Task(name="go_front", parameters=dict(distance=1)),
        Task(name="go_back", parameters=dict(distance=1.55)),
        Task(name="turn_left", parameters=dict(angle=30.0)),
    ]
