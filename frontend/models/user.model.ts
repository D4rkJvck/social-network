export interface User {
    id: number
    email: string
    firstname: string
    lastname: string
    nickname: string
    dateOfBirth: Date
    aboutMe: string
    private: boolean
}

// Id             int
// 	Email          string
// 	Password       string
// 	FirstName      string
// 	LastName       string
// 	Nickname       string
// 	DateOfBirth    time.Time
// 	ProfilePicture string
// 	AboutMe        string
// 	Private        bool