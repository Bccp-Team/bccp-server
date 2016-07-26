Definitions:
=========

Run: execution on a runner
Runner: host accepting a run
Batch: Group of run
namespace: group of repositories

RUNNER_STATE:
 - Waiting
 - Running
 - Dead

RUN_STATE:
 - Waiting (for slots)
 - Running
 - Canceled (by user)
 - Finished (everyting OK)
 - Failed (see return code)
 - Timeout


CI API
======


Runners:
------------

For users:
 - GET /runner -> liste les runners
    return: [ { runner_id: INT, current_run: run_id }, ... ]
  - GET /runner/{id} -> get info sur un runner
    return: { runner_ip: IP, current_run: run_id, runner_run: [ run_id1, ...] }
  - DELETE /runners/{id} -> delete un runner
  - POST /runner/{id}/enable -> reactive un runner
  - POST /runner/{id}/disable -> pause un runner

For runners:
  - PUT /runner/register
    parameters: 
      - CI register token
      - Runner API token
    return: Runner CI API token


Run:
------

For users:
  - GET /run: (lourd)
    return: [ { host: runner_id runner_state: STATE, repo: REPO, logs: INT}, ... ]
  - GET /run/{id}:
    return: [ { host: runner_id, runner_state: STATE, repo: REPO, logs: INT}, ... ]
  - PUT /run: Run given repo with given run.sh
    parameters:
     - repo: [ "repo1", "repo2" ] / At least one
     - repo: namespace_name \ must be defined
     - run.sh: "#!"
    return: batch_id
  - DELETE /run/{id}

For runners:
  - POST /run/{id}/update -> update l'etat d'un job (exit, logs)
    parameters: { state: RUN_STATE, logs: LOG_STRING }

Batch:
---------
  - GET /batch/{id}:
    return: { run_waiting: [ run_id1, run_id2, ... ] ,
                  run_running: [ run_id1, run_id2, ... ],
                  run_finished: [ run_id1, run_id2, ... ],
                  run_failed: [ run_id1, run_id2, ... ] }
  - DELETE /batch/{id}

Namespace:
----------------
  - GET /namespace:
    return: [ { namespace_name: NAME, repositories:[ repo1, repo2, ... ] }, ... ]
  - GET /namespace/{name}:
    return: [ { namespace_name: NAME, repositories:[ repo1, repo2, ... ] }, ... ]
  - PUT /namespace:
    params: { namespace_name: NAME, repositories:[ repo1, repo2, ... ] },
    return: Success|Failure
  - DELETE /namespace/{name}:
    return: [ { namespace_name: NAME, repositories:[ repo1, repo2, ... ] }, ... ]


Runner API:
==========

  - GET /kill -> kill current run
  - GET /ping
  - PUT /run: Launch a run
    parameters:
      - Init
      - Repo
      - RunId
      - UpdateTime
      - Timeout
      - Name


Database:
Runner:
Runner_id | Runner_API_token | Runner_CI_token

Run:
Run_id | State | Runner | Repo | Logs

Batch:
Batch_id | Run_list | namespace_id

Namespace:
namespace_id | namespace_name | repos
