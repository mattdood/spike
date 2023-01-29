<div align="center">
    <img src="https://github.com/mattdood/spike/raw/master/assets/spike.gif" alt="Gif of Spike Spiegel from Cowboy Bebop in space"/>
</div>

**Note:** I do not make any claims to the [Cowboy Bebop](https://en.wikipedia.org/wiki/Cowboy_Bebop) assets, names, or trademarks.

# Spike
A command line TODO application to track my various tasks across systems. The
tasks are stored in a JSON file to ensure portability.

## Installation
Ensure that you have a Golang version `>= 1.19`.

```bash
make install
```

## Usage
The CLI adds, lists, and changes status of tasks. Use the following:

```bash
spike --help

# Create a task
spike create -name "Something" -desc "Longer description"

# Find a task to update (list open tasks)
spike list -status O

# Update a task with ID 1 to closed
spike update -id 1 -status C

# List closed tasks
spike list -status C
```

## Task structure
Tasks are meant to be in either `O` (open) or `C` (closed) status. The structure
is defined below.

**Note:** `id`, `created`, and `updated` are auto generated.

Example `~/spike/tasks.json` structure:
```json
{
    "O": [
        {
            "id": 1,
            "name": "Sample name",
            "description": "Longer description here",
            "created": "2022-01-01",
            "updated": "2022-01-01"
        }
    ],
    "C": [
        {
            "id": 0,
            "name": "Other name",
            "description": "Longer description here",
            "created": "2021-12-01",
            "updated": "2022-01-01"
        }
    ]
}
```

