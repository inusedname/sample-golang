package dto

import "app/data"

type CourseDto struct {
	ID     uint
	Name   string
	Weight uint
}

func FromCourse(course data.Course) CourseDto {
	return CourseDto{
		ID:     course.ID,
		Name:   course.Name,
		Weight: course.Weight,
	}
}

type ClassDto struct {
	ID           uint
	Name         string
	Instructor   data.User
	Course       *data.Course
	SlotTotal    uint
	SlotEquipped uint
}

func FromClass(class data.Class) ClassDto {
	return ClassDto{
		ID:           class.ID,
		Name:         class.Course.Name,
		Instructor:   class.Instructor,
		SlotTotal:    class.SlotTotal,
		SlotEquipped: class.SlotEquipped,
	}
}
