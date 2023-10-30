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

from typing import Dict, Any, List, Optional
from pydantic import BaseModel, Field


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
