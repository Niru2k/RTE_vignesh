## General view about the Todo list  Project

This is a todo list project built using the fiber web framework in Golang. 

The project provides RESTful APIs for user authorization posting tasks 

The backend includes functionalities to sign up, login, task posting, update task , delete task

## Technology Stack

The Job Site Real Time Exercise Task project is built using the following technologies:

- **Golang**: The backend is written in Go (Golang), a statically typed, compiled language.

- **fiber**: The fiber web framework is used to create RESTful APIs and handle HTTP requests.

- **JWT**: JSON Web Tokens are used for secure user authentication and authorization.

- **bcrypt**: Passwords are stored securely in hashed form using the bcrypt hashing algorithm.

- **postgres**: PostgreSQL is an advanced, enterprise class open source relational database that supports both SQL and JSON  querying. 
                It is a highly stable database management system, which has contributed to its high levels of resilience,and correctness. 
   

## Setup

The application will be accessible at `http://localhost:8010`.

 ##  Project explanation

-> the user can sign up and login

-> user can post his task details

-> user can update and delete his task details

-> user can filter date and see his task details


The following API endpoints are available in the Job Site project:

- **POST /signup**: Register a new user account.

- **POST /login**: Log in with registered credentials and receive a JWT token.

- **POST /posttask**: Post task details (user Authorization required).

- **GET /updatetask**: user can update his task details (user Authorization required).

- **GET /deletetask/:id**: user can delete his task details (user Authorization required).

- **PUT /gettasksbydate/:id**: user can use filter date and see his task details(Admin Authentication required).

