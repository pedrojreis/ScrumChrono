name: Bug report 🐛
description: Create a report to help us improve
title: '🐛 '
labels: 'bug'
assignees: ''
body:
    - type: textarea
        validations:
        required: true
        attributes:
        label: Describe the bug
        description: >-
            A clear and concise description of what the bug is.

    - type: textarea
        validations:
        required: true
        attributes:
        label: To Reproduce
        placeholder: >-
                1. Go to '...'
                2. Click on '....'
                3. Scroll down to '....'
                4. See error
        description: >-
            Steps to reproduce the behavior.

    - type: input
        id: scrumchrono_version
        validations:
            required: true
        attributes:
        label: What ScrumChrono version are you using
        description: >
            [e.g. 0.1]

    - type: textarea
        validations:
        required: false
        attributes:
        label: Expected behavior
        description: >-
            A clear and concise description of what you expected to happen.

    - type: input
        id: os
        attributes:
        label: What OS and version are you using
        description: >
            [e.g. macOS Sonoma, Windows 11, Ubuntu 23.10]

    - type: input
        id: terminal
        attributes:
        label: What terminal are you using
        description: >
            [e.g. PowerShell, Windows Terminal, iTerm, Warp]

    - type: textarea
        validations:
        required: false
        attributes:
        label: Additional context
        description: >-
            Add any other context about the problem here.
