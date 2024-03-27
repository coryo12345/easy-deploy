# easy-deploy
Easy tool to set up ci/cd on a server with docker. This project was built to simplify my own life, not necessarily meant to be used by others, but it should work for anyone.

NOTE: This is not a tool for production sites. It does not make backup images for previous deploys. 
This is intended to be used only for hobbyists who want to self-host their projects. 

Only consider using this tool if:
* Your code repositories on stored on GitHub and you can set up an action to trigger a webhook
* You do not need to store historic images for projects
* You want to run each project in docker container(s)
    * you may split up a project into multiple services with different dockerfiles
    * but your application code must support running & communicating between containers

You will still need to configure any reverse proxies yourself to point traffic to each container (if needed)

## Info
The example_config.json example has all properties that are used defined. It should be pretty straightforward, the details provided for each service will be used to follow a basic CLONE -> BUILD -> STOP_OLD -> DELETE_OLD -> START_NEW series of steps.

## How to Run
Duplicate the `.env_example` file to a `.env` file and adjust settings as needed. Ensure the `DEPLOY_ENV_ENVIRONMENT` variable is set to `prod`.  
Adjust the `example_config.json` or create your own.  
Install the needed dependencies with `make install`.  
Run the code with `make run` or `make build && ./main`.  


## TODO
* [ ] need to make `start.sh` file that runs easy-deploy in a docker container (must mount docker socket & config file)
* [ ] research github webhooks to trigger deploys
