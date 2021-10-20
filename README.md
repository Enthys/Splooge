# Wildfire
## About
**Wildfire** is a mass update tool. You can register projects from git/gitlab, add those projects to groups then through
**Wildfire** clone in parallel those projects and execute commands on all(or only some) projects at once. This tool is
created to solve the issue with manually going through all services and updating them one by one, with **Wildfire** you
can update all the services/packages/tools at once.

---

## Configuration
A configuration is composed of 2 sets one for `Projects` and one for `Groups`
Example configuration:
```yaml
groups:
  node_apps:
     - foo
     - bar
  http_services:
     - foo
     - zaz

projects:
  foo:
    name: foo
    type: git
    url: git@github.com/example/foo
  bar:
    name: bar
    type: gitlab
    url: git@gitlab.com/example/bar
  zaz:
    name: zaz
    type: bit
    url: git@bitbucket.com/example/zaz
```

## Usage
### Global Arguments
 - `--config` - Specify which configuration to use. If not set **Wildfire** will create a new configuration in current
directory under the name `.wildfire.yaml`

### Add Project
Will create a new project record in configuration.  
If a project with the same name already exists it will return an 
error
```shell
$ wildfire project add <name> <type> <ssh-address>
```
 #### Parameters
 - `name` - The name of the project.
 - `type` - _(Currently not used but required)_ Source of the project. Available options:
   - `git`
   - `gitlab`
   - `bitbucket`
 - `ssh-address` - The address from which to retrieve the project. Will be used with the `git clone` command

---

### Get Projects
Displays all registered projects in configuration
```shell
$ wildfire projects list
```

---

### Remove Project
Will remove provided projects from the configuration and from groups.
```shell
$ wildfire project remove <name>...
```
#### Parameters
 - `name` - A list of the projects which we want to remove

---

### Update or Create Project
If the project does not exist it will create a new project record.  
If the project does exist, then it will prompt for used input whether to overwrite the project record or not.
```shell
$ wildfire project set <name> <type> <ssh-address>
```

---

### Create Group
```shell
$ wildfire group create <name> [project-name]...
```
#### Parameters
 - `name` - The name of the group we want to create
 - `project-name` - _(optional)_ A list of the projects we want to add to the group after creation. Will cancel group
creation if a project which does not exist is provided.

---

### Delete Group
```shell
$ wildfire group delete <name>
```
#### Parameters
 - `name` - The name of the group we want to delete

---

### Add Project to Group
```shell
$ wildfire group add-project <name> <project-name>...
```
#### Parameters
 - `name` - The name of the group to which we want to add the projects. Will throw an error if we provide group which
does not exist
 - `project-name` - At least one is required. The names of the projects we want ot add to the group. Will not update 
group if a project does not exist in configuration

---

### Remove Project from Group
```shell
$ wildfire group remove-project <name> <project-name>...
```
#### Parameters
 - `name` - The name of the group from which we want to remove projects
 - `project-name` - At least one is required. The names of the projects we want to remove from the group. Will not throw
an error if a project is not found in group or configuration.
