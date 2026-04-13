lines = open("internal/repository/user/get.go", "r").readlines()
with open("internal/repository/user/get.go", "w") as f:
    f.writelines(lines[:105])
