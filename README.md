# notes
Interview coding challenge

###Installation
Set the port in ```config/config.yml```
```
go build start.go
./start
```

###API
| Function           | HTTP Method | Path              | Response / Reqest                                                                              |
|--------------------|-------------|-------------------|------------------------------------------------------------------------------------------------|
| New User           | `GET`       | `/users/new`      | **RESPONSE** *userId* <br/> 36 character guid belonging to a newly created user on the server. |
| New Note           | `GET`       | `/notes/new`      | **RESPONSE** *noteId* <br/> 36 character guid belonging to a newly created note on the server. |
| Update Note        | `PUT`       | `/notes`          | **REQUEST** NoteMessage <br/>                                                                  |
| Get Saved Notes    | `GET`       | `/notes/saved`    | **REQUEST** User<br/>**RESPONSE** NoteCollection                                               |
| Get Archived Notes | `GET`       | `/notes/archived` | **REQUEST** User<br/>**RESPONSE** NoteCollection                                               |

###Objects

**User**
```json
{
  "userId": "10bfa4b4-0d40-4e2c-b458-548cd11c85de", 
  "userName": "John Smith"
}
```
**Note**
```json
{
  "noteId": "2b9a81b0-5dc8-4022-b06d-32d33dc06268",
  "title": "New Note",
  "content": "Shopping List\nCheese",
  "archived": false
}
```
**NoteMessage**
```json
{
  "user": User,
  "note": Note
}
```
**NoteCollection**
```json
{
  "notes": Note[]
}
```