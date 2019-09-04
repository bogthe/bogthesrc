# bogthesrc 

A news website built with Go.

## Features
Here come a list of features, buzz-words to get excited about:
- Open Source Software; 
- Built with Go 1.12: modules + testing (a pragmatic approach);
- RESTful JSON Api to serve posts, serve a post by ID and creation of posts;
- Deterministic development environment aka Docker with docker-compose to create a development environment which uses reflex for reloading of code and spins up two services: app and PostgreSQL; 
- WebApp built using go templates with CSS3 Flexbox and Grids;
- A Command Line Interface (cli) to spin up the server, create / drop the database,create a basic post, start the importer service;
- Importers: asynchronous workers that scrape the hacker news API, transforms the data, and inserts it into our database;

## Reason
The main reason for this website is to have fun and build something useful using the beautiful language / ecosystem that Go provides. The main source of inspiration is a project called `sourcegraph/thesrc`. 

I went through every single one of the project's commits from commit 1 to final to learn how their project evolved. It taught me a lot about what a production app might be structured around: folder structures, code structure, abstractions, DB interaction etc. Some abstractions I kept, some I discarded and other functionality was added specifically importing from the hackernews API endpoints. 

I will be sharing all my learnings on my personal blog [bogthe.me]("https://bogthe.me/")

![Alt text](./bogthesrc.png?raw=true "Title")
