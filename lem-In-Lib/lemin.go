package lemIn

import (
	"path"
)

func IsFileTypeCorrect(givenFileName string) bool {
	extension := path.Ext(givenFileName)
	return extension == ".txt"
}

func CleanTheComments(seperatedContent []string) []string {
	var res []string
	for i := 0; i < len(seperatedContent); i++ {
		if seperatedContent[i][0] == '#' && seperatedContent[i][1] != '#' {
			continue
		} else {
			res = append(res, seperatedContent[i])
		}
	}
	return res
}

func SeperateTheContent(content string) []string {
	var res []string
	var line string
	for _, char := range content {
		if char == '\n' {
			res = append(res, line)
			line = ""
		} else {
			line += string(char)
		}
	}
	if line != "" {
		res = append(res, line)
	}
	return res
}

func IsFormatOk(seperatedContent []string) string {
	if AreThereAnyAnt(seperatedContent) != "" {
		return AreThereAnyAnt(seperatedContent)
	}
	if !HaveStartAndEnd(seperatedContent) {
		return "There is problem because of ##start and ##end points!"
	}
	if AreThereMoreOrFewThanTwoDoubleHash(seperatedContent) != "" {
		return AreThereMoreOrFewThanTwoDoubleHash(seperatedContent)
	}
	if AreRoomsDuplicated(seperatedContent) {
		return "There is duplicated rooms!"
	}
	if AreTheCoordinatesValid(seperatedContent) != "" {
		return AreTheCoordinatesValid(seperatedContent)
	}
	if !ControlRoomNamesStart(seperatedContent) {
		return "Room's name cannot start with '#' or 'L'"
	}
	if AreThereSameLink(seperatedContent) {
		return "A room can connect only one time!"
	}
	if IsEndRoomAlone(seperatedContent) {
		return "There is no connection to ##end room!"
	}
	if IsThereAnyUnknownLink(seperatedContent) {
		return "There is unknown room in the links!"
	}

	return ""
}

func AreThereAnyAnt(seperatedContent []string) string {
	if seperatedContent[0] == "0" {
		return "There is no ant!"
	}

	for _, char := range seperatedContent[0] {
		if char < '0' || char > '9' {
			return "Invalid chars in the ant numbers!"
		}
	}

	return ""
}

func HaveStartAndEnd(seperatedContent []string) bool {
	haveStart := false
	haveEnd := false
	for _, line := range seperatedContent {
		if line == "##start" {
			haveStart = true
		} else if line == "##end" {
			haveEnd = true
		}
	}
	return haveStart && haveEnd
}

func IsItRoom(line string) bool {
	spaceCount := 0
	for _, char := range line {
		if char == ' ' {
			spaceCount++
		}
	}
	return spaceCount == 2
}

func IsItLink(line string) bool {
	for i := 0; i < len(line); i++ {
		if line[i] == '-' {
			if line[i-1] != ' ' && line[i+1] != ' ' {
				return true
			}
		}
	}
	return false
}

func TakeTheStartAndEndRooms(seperatedContent []string) [2]string {
	var res [2]string
	for i := 0; i < len(seperatedContent); i++ {
		if seperatedContent[i] == "##start" {
			res[0] = seperatedContent[i+1]
		} else if seperatedContent[i] == "##end" {
			res[1] = seperatedContent[i+1]
		}
	}
	return res
}

func TakeTheLinks(seperatedContent []string) []string {
	var res []string
	for _, line := range seperatedContent {
		if IsItLink(line) {
			res = append(res, line)
		}
	}
	return res
}

func ReverseTheString(link string) string {
	res := ""
	for i := len(link) - 1; i >= 0; i-- {
		res += string(link[i])
	}
	return res
}

func TakeTheRoomName(line string) string {
	res := ""
	for _, char := range line {
		if char != ' ' {
			res += string(char)
		} else {
			break
		}
	}
	return res
}

func AreRoomsDuplicated(seperatedContent []string) bool {
	var roomNumbers []string
	for _, line := range seperatedContent {
		if IsItRoom(line) {
			roomNumbers = append(roomNumbers, TakeTheRoomName(line))
		}
	}

	for i := 0; i < len(roomNumbers)-1; i++ {
		for j := i + 1; j < len(roomNumbers); j++ {
			if roomNumbers[i] == roomNumbers[j] {
				return true
			}
		}
	}
	return false
}

func AreTheCoordinatesValid(seperatedContent []string) string {
	// two rooms can't be at the same coordinate
	var roomsCoordinate []string

	for _, line := range seperatedContent {
		tempRes := ""
		addToList := false
		if IsItRoom(line) {
			spaceCount := 0
			for i := len(line) - 1; i >= 0; i-- {
				if string(line[i]) == " " {
					spaceCount++
				}
				if spaceCount == 2 {
					addToList = true
					break
				}
				tempRes = string(line[i]) + tempRes
			}
		}
		if addToList {
			roomsCoordinate = append(roomsCoordinate, tempRes)
		}
	}

	for i := 0; i < len(roomsCoordinate)-1; i++ {
		for j := i + 1; j < len(roomsCoordinate); j++ {
			if roomsCoordinate[i] == roomsCoordinate[j] {
				return "Rooms cannot be at the same coordinates!"
			}
		}
	}

	// coordinates can't be negative and can't be chars other than nums
	for _, coordinate := range roomsCoordinate {
		for _, char := range coordinate {
			if (char >= '0' && char <= '9') || char == ' ' {
				continue
			} else {
				return "Invalid coordinates!"
			}
		}
	}

	return ""
}

func CreateLinkWithStartAndEndRooms(startAndEndRooms [2]string) []string {
	var res []string
	res = append(res, TakeTheRoomName(startAndEndRooms[0])+"-"+TakeTheRoomName(startAndEndRooms[1]))
	res = append(res, TakeTheRoomName(startAndEndRooms[1])+"-"+TakeTheRoomName(startAndEndRooms[0]))
	return res
}

func ControlRoomNamesStart(seperatedContent []string) bool {
	var roomNames []string

	for _, line := range seperatedContent {
		if IsItRoom(line) {
			roomNames = append(roomNames, TakeTheRoomName(line))
		}
	}

	for _, name := range roomNames {
		if name[0] == '#' || name[0] == 'L' {
			return false
		}
	}
	return true
}

func AreThereSameLink(seperatedContent []string) bool {
	links := TakeTheLinks(seperatedContent)

	for i := 0; i < len(links)-1; i++ {
		revStr := ReverseTheString(links[i])
		for j := i + 1; j < len(links); j++ {
			if revStr == links[j] {
				return true
			}
		}
	}
	return false
}

func IsThereAnyUnknownLink(seperatedContent []string) bool {
	var roomsNames []string
	links := TakeTheLinks(seperatedContent)

	for _, line := range seperatedContent {
		if IsItRoom(line) {
			roomsNames = append(roomsNames, TakeTheRoomName(line))
		}
	}

	var namesFromLinks []string

	for _, line := range links {
		tempRoom := ""
		for _, char := range line {
			if char != '-' {
				tempRoom += string(char)
			} else {
				namesFromLinks = append(namesFromLinks, tempRoom)
				tempRoom = ""
			}
		}
		if tempRoom != "" {
			namesFromLinks = append(namesFromLinks, tempRoom)
		}
	}

	for _, nameFromLink := range namesFromLinks {
		thereIsUnknown := true
		for _, roomName := range roomsNames {
			if nameFromLink == roomName {
				thereIsUnknown = false
				break
			}
		}
		if thereIsUnknown {
			return true
		}
	}

	for _, roomName := range roomsNames {
		thereIsUnknown := true
		for _, nameFromLink := range namesFromLinks {
			if nameFromLink == roomName {
				thereIsUnknown = false
				break
			}
		}
		if thereIsUnknown {
			return true
		}
	}
	return false
}

func IsEndRoomAlone(seperatedContent []string) bool {
	startAndEndRooms := TakeTheStartAndEndRooms(seperatedContent)
	endRoomName := TakeTheRoomName(startAndEndRooms[1])

	var roomsInLinks []string

	for _, line := range seperatedContent {
		tempRoomName := ""
		if IsItLink(line) {
			for _, char := range line {
				if char != '-' {
					tempRoomName += string(char)
				} else {
					roomsInLinks = append(roomsInLinks, tempRoomName)
					tempRoomName = ""
				}
			}
			if tempRoomName != "" {
				roomsInLinks = append(roomsInLinks, tempRoomName)
			}
		}
	}

	for _, room := range roomsInLinks {
		if room == endRoomName {
			return false
		}
	}

	return true
}

func AreThereMoreOrFewThanTwoDoubleHash(seperatedContent []string) string {
	hashLineCount := 0

	for _, line := range seperatedContent {
		if line[0] == '#' && line[1] == '#' {
			hashLineCount++
		}
	}

	if hashLineCount < 2 {
		return "There are double hashes less than two!"
	} else if hashLineCount > 2 {
		return "There are double hashes more than two!"
	}
	return ""
}

func ClearTheLinks(seperatedContent []string) []string {
	var res []string

	for i := 0; i < len(seperatedContent); i++ {
		linkExist := false
		if i == len(seperatedContent)-1 {
			res = append(res, seperatedContent[i])
			linkExist = true
		} else if IsItLink(seperatedContent[i]) {
			for j := i + 1; j < len(seperatedContent); j++ {
				if seperatedContent[i] == seperatedContent[j] {
					linkExist = true
				}
			}
		} else {
			res = append(res, seperatedContent[i])
			linkExist = true
		}

		if !linkExist {
			res = append(res, seperatedContent[i])
		}
	}
	return res
}
