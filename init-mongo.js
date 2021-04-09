db.createUser(
    {
        user: "articles",
        pwd: "articles_pass",
        roles: [
            {
                role: "readWrite",
                db: "articles"
            }
        ]
    }
)