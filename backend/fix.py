with open("internal/repository/user/get.go", "a") as f:
    f.write("""

func (r *UserRepo) GetByIDs(ids []int64) ([]user.User, error) {
    if len(ids) == 0 {
        return []user.User{}, nil
    }

    placeholders := ""
    args := make([]interface{}, len(ids))
    for i, id := range ids {
        if i > 0 {
            placeholders += ","
        }
        placeholders += "?"
        args[i] = id
    }

    query := `SELECT id, email, first_name, last_name, date_of_birth, avatar_url, nickname, about_me, is_private FROM users WHERE id IN (` + placeholders + `)`

    rows, err := r.db.Query(query, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var users []user.User
    for rows.Next() {
        var u user.User
        if err := rows.Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.DOB, &u.AvatarURL, &u.Nickname, &u.AboutMe, &u.IsPrivate); err != nil {
            return nil, err
        }
        users = append(users, u)
    }

    if err = rows.Err(); err != nil {
        return nil, err
    }

    return users, nil
}
""")
