db.createUser(
    {
        user: "user_articles",
        pwd: "user_articles_pass",
        roles: [
            {
                role: "readWrite",
                db: "user_articles"
            }
        ]
    }
)