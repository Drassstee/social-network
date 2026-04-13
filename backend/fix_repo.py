import re

content = open("internal/repository/repository.go", "r").read()

content = content.replace('"social-network/internal/repository/user"', 'userrepo "social-network/internal/repository/user"')
content = content.replace('"social-network/internal/repository/post"', 'postrepo "social-network/internal/repository/post"')
content = content.replace('user.NewUserRepo', 'userrepo.NewUserRepo')
content = content.replace('post.NewPostRepo', 'postrepo.NewPostRepo')

with open("internal/repository/repository.go", "w") as f:
    f.write(content)
