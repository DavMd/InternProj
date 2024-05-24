# InternProj

InternProj is a GraphQL-based application that allows users to create posts and comments.
It is built using Go and uses PostgreSQL for the database.
This project is containerized using Docker for easy setup and deployment.

## Features

- Create, read, update, and delete posts
- Comment on posts
- Nested comments with support for threading
- Disable comments on posts
- Subscribing on post comments

## Technologies

- **Backend:** Go
- **Database:** PostgreSQL
- **API:** GraphQL
- **Containerization:** Docker, Docker Compose

## Getting Started

### Prerequisites

- Docker
- Docker Compose

### Clone the repository

```sh
git clone https://github.com/yourusername/InternProj.git
cd InternProj
```

### Build and run the Docker containers

```sh
docker-compose up --build
```

This command will build the Docker images and start the containers.
The application will be available at http://localhost:8080 and PostgreSQL will be available at localhost:5432.

# GraphQL API

## Queries

### Get all posts

```graphql
query {
  posts {
    id
    title
    body
    is_disabled_comments
    user_id
  }
}
```

### Get a post by ID

```graphql
query($id: UUID!) {
  post(id: $id) {
    id
    title
    body
    is_disabled_comments
    user_id
  }
}
```

## Mutations

### Create a post

```graphql
mutation($title: String!, $body: String!, $user_id: String!) {
  createPost(title: $title, body: $body, user_id: $user_id) {
    id
    title
    body
    is_disabled_comments
    user_id
  }
}
```

### Create a comment

```graphql
mutation($post_id: UUID!, $body: String!, $user_id: String!, $parent_id: UUID) {
  createComment(post_id: $post_id, body: $body, user_id: $user_id, parent_id: $parent_id) {
    id
    post_id
    parent_id
    body
    user_id
  }
}
```

## Subscriptions

### Comment added to a post

```graphql
subscription($postID: ID!) {
  commentAdded(postID: $postID) {
    id
    post_id
    parent_id
    body
    user_id
  }
}
```

# Database Schema

## Posts Table

```sql
CREATE TABLE posts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    body TEXT NOT NULL,
    is_disabled_comments BOOLEAN DEFAULT FALSE,
    user_id TEXT NOT NULL
);
```

## Comments Table

```sql
CREATE TABLE comments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    post_id UUID NOT NULL,
    parent_id UUID,
    body TEXT NOT NULL,
    user_id TEXT NOT NULL,
    CONSTRAINT fk_post
      FOREIGN KEY(post_id) 
      REFERENCES posts(id)
      ON DELETE CASCADE,
    CONSTRAINT fk_parent_comment
      FOREIGN KEY(parent_id)
      REFERENCES comments(id)
      ON DELETE CASCADE
);

CREATE INDEX idx_post_id ON comments(post_id);
CREATE INDEX idx_parent_id ON comments(parent_id);
```
