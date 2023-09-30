from setuptools import setup
from glob import glob
import os

package_name = "bot_description"

setup(
    name=package_name,
    version="0.0.0",
    packages=[package_name],
    data_files=[
        ("share/ament_index/resource_index/packages", ["resource/" + package_name]),
        ("share/" + package_name, ["package.xml"]),
        (os.path.join("share", package_name, "launch"), glob("launch/*.launch.py")),
        (os.path.join("share", package_name, "urdf"), glob("urdf/**")),
        (os.path.join("share", package_name, "world"), glob("world/**")),
    ],
    install_requires=["setuptools", "promptulate"],
    zip_safe=True,
    description="TODO: Package description",
    license="TODO: License declaration",
    tests_require=["pytest"],
    entry_points={
        "console_scripts": [
            "rotate_wheel= bot_description.rotate_wheel:main",
            "user_client= bot_description.user_client:main",
            "odom_client= bot_description.odom_client:main",
        ],
    },
)
