package user

func (r *UserRepo) UserExists(id int64) (bool, error) {
	query := `SELECT EXISTS(
				SELECT 1
				FROM users
				WHERE id = ?)`

	var exists int
	err := r.db.QueryRow(query, id).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists == 1, nil
}

// --------------------------------------------------------------------|

func (r *UserRepo) EmailExists(email string, id int64) (bool, error) {
	query := `SELECT EXISTS(
				SELECT 1
				FROM users
				WHERE email = ? AND id != ?)`

	var exists int
	err := r.db.QueryRow(query, email, id).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists == 1, nil
}

// --------------------------------------------------------------------|

func (r *UserRepo) IsPrivate(id int64) (bool, error) {
	query := `SELECT EXISTS(
				SELECT 1 
				FROM users 
				WHERE id = ? AND profile_type = 'private')`

	var exists int
	err := r.db.QueryRow(query, id).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists == 1, nil
}
