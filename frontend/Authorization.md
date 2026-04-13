# End Points

Register - `localhost:8080/api/v1/register`
Login - `localhost:8080/api/v1/login`
Logout - `localhost:8080/api/v1/logout`
Delete - `localhost:8080/api/v1/delete`
Get Profile - `localhost:8080/api/v1/users/{id}`
Update Profile - `localhost:8080/api/v1/users`
Notifications - `localhost:8080/api/v1/notifications`
Follow - `localhost:8080/api/v1/follow`
Unfollow - `localhost:8080/api/v1/unfollow`
UploadAvatar - `localhost:8080/api/v1/avatar`
GetAvatar - `localhost:8080/api/v1/avatar/{id}`

---
HTTP Status Code: 200, 201, 400, 403, 404, 409

Avatar отправлят через `multipart/form-data`. Загрузка аватара всегда отдельно происходит. `Notifications` можно периодический вызывать для  private профилю, чтобы занать есть ли уведоление о запросе на подписку.

---
## Register
URL : `localhost:8080/api/v1/register`
**Method** : `POST`
**Request** : 
```json
{
	"email": "mkapan@gmail.com",
	"first_name": "mkapan",
	"last_name": "mkapan",
	"password": "Qwerty",
	"dob": "2026-03-31T15:04:47.645491103+05:00",
	"avatar": "", // может быть пустым
	"nickname": "", // может быть пустым
	"about_me": "" // может быть пустым
}
```

**Response** : 
**Status** : 201 
```json
{
	"id": 1,
	"first_name": "mkapan",
	"last_name": "mkapan",
	"expires_at": "2026-04-01T16:40:08.614763942+05:00"
}
```

---
## Login
URL : `localhost:8080/api/v1/login`
Method : `POST`
Request :
```json
{
	"email": "mkapan@gmail.com",
	"password": "Qwerty"
}
```
Response :
Status : 200
```json
{
	"id": 1,
	"first_name": "mkapan",
	"last_name": "mkapan",
	"expires_at": "2026-04-01T17:07:24.70943009+05:00"
}
```

---
## Logout
URL : `localhost:8080/api/v1/logout`
Method : `POST`
Request :
```json
// Empty
```
Response :
Status : 204
```json
// null
```

---
## Delete
URL : `localhost:8080/api/v1/delete`
Method : `DELETE`
Request :
```json
// Empty
```
Response :
Status : 204
```json
// null
```

---
## GetProfile
URL : `localhost:8080/api/v1/users/{id}`
Method : `GET`
Request :
```json
// Empty
```
Response :
```json
// Если профиль private и ты не подписчик его
{
	"user": {
		"id": 3,
		"first_name": "test2",
		"last_name": "test2",
		"profile_type": "private"
	},
	"followers": null,
	"following": null,
	"posts": null
}

// Когда профиль public или private, но ты его подписчик
{
	"user": {
		"id": 1,
		"email": "mkapan@gmail.com",
		"first_name": "mkapan",
		"last_name": "mkapan",
		"dob": "2026-03-31T15:04:47.645491103+05:00",
		"profile_type": "public"
	},
	"followers": [
		{
			"id": 2,
			"first_name": "test",
			"last_name": "test1"
		},
		{
			"id": 3,
			"first_name": "mkapan",
			"last_name": "mkapan"
		}
	],
	"following": [
		{
			"id": 3,
			"first_name": "mkapan",
			"last_name": "mkapan"
		}
	],
	"posts": null
}
```

---
## UpdateProfile
URL : `localhost:8080/api/v1/users`
Method : `PUT`
Request :
```json
{
	"email": "mkapan@gmail.com",
	"first_name": "mkapan",
	"last_name": "mkapan",
	"dob": "2026-03-31T15:04:47.645491103+05:00",
	"avatar": "",
	"nickname": "MK",
	"about_me": "",
	"profile_type": "private"
}
```
Response :
```json
{
	"id": 1,
	"email": "mkapan@gmail.com",
	"first_name": "mkapan",
	"last_name": "mkapan",
	"dob": "2026-03-31T15:04:47.645491103+05:00",
	"avatar": "",
	"nickname": "MK",
	"about_me": "",
	"profile_type": "private"
}
```

---
## Follow
Если followers_id = user.ID значит я его подписчик.

URL : `localhost:8080/api/v1/follow`
Method : `POST`
Request :
```json
{
	"following_id": 1, // на кого я подписываюсь
}
```
Response :
Status : 200
```json
{
	"status": "accept"
}
```

---
## Unfollow
URL : `localhost:8080/api/v1/unfollow`
Method : `POST`
Request :
```json
{
	"following_id": 1 // От кого я отписываюсь
}
```
Response :
Status : 200
```json
{
	"status": "unfollow"
}
```

---
## GetNotifications
URL : `localhost:8080/api/v1/notifications
Method : GET
Request :
```json
// Empty
```
Response :
```json
// Если есть люди, которые хотят подписаться. Только если профиль private
[
	{
		"id": 1,
		"first_name": "mkapan",
		"last_name": "mkapan"
	}
]
```

---
## ResponseNotifications
URL : `localhost:8080/api/v1/notifications
Method : POST
Request :
```json
// Accept/Decline
{
	"follower_id": 1,
	"status": "accept"
}
```
Response :
```json
// Empty
```

---
## UploadAvatar
URL : `localhost:8080/api/v1/avatar
Method : POST
Request : через `multipart/form-data`
Response :
```json
// Empty
```

---
## GetAvatar
URL : `localhost:8080/api/v1/avatar/{id}
Method : GET
Request :
```json
// Empty
```
Response :
```json
// Отправится фото
```

---


# Error 
## Register
URL : `localhost:8080/api/v1/register

Response :
Status : 409
```json
// Когда email уже существует. Email должен быть уникальным.
{
	"error": "register: conflict: email already exists"
}
```
Response :
Status : 400
```json
// Когда email не правильный
{
	"error": "register: invalid data: incorrect email"
}

// Когда first_name пустой
{
	"error": "register: invalid data: first name is empty"
}

// Когда last_name пустой
{
	"error": "register: invalid data: last name is empty"
}

// Когда password пустой
{
	"error": "register: invalid data: password is empty"
}

// Когда dob(date_of_birth) пустой или не правильный
{
	"error": "parsing time \"\" as \"2006-01-02T15:04:05Z07:00\": cannot parse \"\" as \"2006\"" // Что-то похожое ошибка будет
}
```

---
## Login
URL : `localhost:8080/api/v1/login`

Response :
Status : 400
```json
// Когда не существует email или не правильный пароль password
{
	"error": "login: invalid data: invalid email or password"
}

// Когда не правильный email 
{
	"error": "login: invalid data: incorrect email"
}

// Когда email пустой
{
	"error": "login: invalid data: email is empty"
}

// Когда password пустой
{
	"error": "login: invalid data: password is empty"
}
```

---
## Logout
URL : `localhost:8080/api/v1/logout`

Response :
Status : 401
```json
// Если после выхода ещё раз выйти, то есть не авторизован пользователь
{
	"error": "no cookie"
}
```

---
## Follow
URL : `localhost:8080/api/v1/follow`

Response :
Status : 400
```json
// Когда попытался подписатся на себя
{
	"error": "follow: invalid data: self-following is not allowed"
}

// Когда id не правильный
{
	"error": "follow: invalid data: incorrect user id"
}
```
Status : 404
```json
// Когда пользователя не существует
{
	"error": "follow: not found: user not found"
}
```

Status : 409
Response: 
```json
// Когда уже существует подписка
{
	"error": "follow: conflict: follow already exists"
}
```


---
## Unfollow
URL : `localhost:8080/api/v1/unfollow

Response :
Status : 400
```json
// Когда id не правильный или не cуществует
{
	"error": "follow: invalid data: incorrect user id"
}

// Когда не cуществует подписки
{
	"error": "unfollow: invalid data: not following this user"
}

// Когда попытался отписаться от себя
{
	"error": "follow: invalid data: self-unfollowing is not allowed"
}
```

---
## GetProfile
URL : `localhost:8080/api/v1/users/{id}`

Status: 404
Response :
```json
// Когда id не правильный
{
	"error": "get profile: not found: user not found"
}
```

---
## UpdateProfile
URL : `localhost:8080/api/v1/users`

Status : 409
Response :
```json
// Когда email уже существует. Email должен быть уникальным.
{
	"error": "register: conflict: email already exists"
}
```
Status : 400
Response :
```json
// Когда email не правильный
{
	"error": "register: invalid data: incorrect email"
}

// Когда first_name пустой
{
	"error": "register: invalid data: first name is empty"
}

// Когда last_name пустой
{
	"error": "register: invalid data: last name is empty"
}

// Когда dob(date_of_birth) пустой или не правильный
{
	"error": "parsing time \"\" as \"2006-01-02T15:04:05Z07:00\": cannot parse \"\" as \"2006\"" // Что-то похожое ошибка будет
}

// Когда профиль тип не public и не private
{
	"error": "update profile: invalid data: incorrect profile type"
}
```

---
## ResponseNotifications
URL : `localhost:8080/api/v1/notifications`

Status : 400
Response :
```json
// Когда id не правильный
{
	"error": "notification: invalid data: incorrect user id"
}

// Когда status не правильный
{
	"error": "notification: invalid data: incorrect status"
}
```
Status : 404
```json
// Когда не существует такого записи в БД
{
	"error": "notification: not found: follow request not found"
}
```

---
## UploadAvatar
URL : `localhost:8080/api/v1/avatar

Response :
```json
// Empty
```

---
## GetAvatar
URL : `localhost:8080/api/v1/avatar/{id}

Status : 400
Response :
```json
// Когда id не правильный
{
	"error": "strconv.ParseInt: parsing \"ф\": invalid syntax"
}

```
Status : 404
```json
// Когда пользователя не существует
{
	"error": "get avatar: not found: user not found"
}
```

---
