# easy-deploy
Easy tool to set up ci/cd on a server with docker

NOTE: This is not a tool for production sites. It does not make backup images for previous deploys. 
This is intended to be used only for hobbyists who want to self-host their projects. 

Only consider using this tool if:
* Your code repositories on stored on GitHub and you can set up an action to trigger a webhook
* You do not need to store historic images for projects
* You want to run each project in docker container(s)
    * you may split up a project into multiple services with different dockerfiles
    * but your application code must support running & communicating between containers

You will still need to configure any reverse proxies yourself to point traffic to each container (if needed)

## TODO
* [ ] need to make `start.sh` file that runs easy-deploy in a docker container (must mount docker socket & config file)
* [ ] research github webhooks to trigger deploys
