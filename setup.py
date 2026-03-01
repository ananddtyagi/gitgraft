from setuptools import setup, find_packages

with open("README.md", "r", encoding="utf-8") as fh:
    long_description = fh.read()

setup(
    name="gitgraft",
    version="0.1.0",
    author="GitGraft Contributors",
    description="TUI for git commands",
    long_description=long_description,
    long_description_content_type="text/markdown",
    url="https://github.com/ananddtyagi/gitgraft",
    packages=find_packages(),
    classifiers=[
        "Programming Language :: Python :: 3",
        "License :: OSI Approved :: MIT License",
        "Operating System :: OS Independent",
    ],
    python_requires=">=3.8",
    install_requires=[
        "textual>=0.47.0",
        "GitPython>=3.1.40",
    ],
    entry_points={
        "console_scripts": [
            "git-graft=gitgraft.cli:main",
        ],
    },
)
