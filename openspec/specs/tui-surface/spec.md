# TUI Surface Specification

## Purpose

Define the interactive TUI Surface behavior for task-first, people-facing workflows.

## Requirements

### Requirement: Initial TUI shell shows profile and primary actions

The TUI Surface MUST start from the shared Application and render a task-first shell that shows the effective Profile and curated primary Actions.

#### Scenario: People-facing TUI-visible capability is shown

- GIVEN a Capability has people audience and TUI visibility
- WHEN the initial TUI model is rendered
- THEN that Capability appears as a primary Action

#### Scenario: Agent-only capability is not promoted as primary

- GIVEN a Capability has only agent audience
- WHEN the initial TUI model is rendered
- THEN that Capability is not shown as a primary Action
