package model

type User struct {
    ID        string `json:"id"`
    Email     string `json:"email"`
    Password  string `json:"password"`
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
    DOB       string `json:"dob"`
    Avatar    string `json:"avatar,omitempty"`
    Nickname  string `json:"nickname,omitempty"`
    AboutMe   string `json:"about_me,omitempty"`
}