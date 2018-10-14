package common

const (
	// Scrapped from Faculty Directory
	tableCreationDepartments = `CREATE TABLE departments (
		id SERIAL PRIMARY KEY,
		code varchar(5) NOT NULL,	-- 2 character code of each department
		name varchar(95) NOT NULL		-- Full name of the department
	);`

	// Scrapped from Faculty Directory
	tableCreationFacultyDesignation = `CREATE TABLE faculty_designations (
		id SERIAL PRIMARY KEY,
		designation varchar(80) NOT NULL	-- Designations in KGP
	);`

	// Pre-populated + Scrapped
	tableCreationRooms = `CREATE TABLE rooms (
		id SERIAL PRIMARY KEY,
		room varchar(80) NOT NULL		-- Room Name/Room No
	);`

	// Pre-populated manually.
	// 9 slots daily over a work week => 9 * 5 = 45 slots
	tableCreationTimeSlots = `CREATE TABLE time_slots (
		id SERIAL PRIMARY KEY,
		day varchar(10) NOT NULL,	-- Week Day
		slot varchar(5) NOT NULL	-- Time slots in a working day e.g 8 AM
	);`

	// Scrapped from Faculty Directory
	// Every faculty member should have a designation and a department.
	tableCreationFaculty = `CREATE TABLE faculty (
		id SERIAL PRIMARY KEY,
		name varchar(80) NOT NULL,
		designation int REFERENCES faculty_designations(id),
		department int REFERENCES departments(id)
	);`

	// Faculty should exist for a course to exist.
	// Multiple rows corresponding to the same course can exist in the table.
	// - Same course with different course codes.
	// - Multiple faculty members teaching the same course.
	tableCreationCourses = `CREATE TABLE courses (
		id SERIAL PRIMARY KEY,
		code varchar(10) NOT NULL,
		name varchar(80) NOT NULL,
		credits int,
		faculty int REFERENCES faculty(id),
		department int REFERENCES departments(id)
	);`

	// Course must exist to show up in the timetable.
	// Every course must have a time slot and an alloted room.
	tableCreationAcademicTimetable = `CREATE TABLE academic_timetable (
		id SERIAL PRIMARY KEY,
		course int REFERENCES courses(id),
		time_slot int REFERENCES time_slots(id),
		room int REFERENCES rooms(id)
	);`
)

const (
	TableInsertionDepartment  = `INSERT INTO departments (code, name) VALUES ($1, $2);`
	TableInsertionDesignation = `INSERT INTO faculty_designations (designation) VALUES ($1);`
	TableInsertionFaculty     = `INSERT INTO faculty (name, designation, department) VALUES ($1, $2, $3);`
	TableInsertionCourses     = `INSERT INTO courses (code, name, credits, faculty, department) VALUES ($1, $2, $3, $4, $5);`
	TableInsertionRooms       = `INSERT INTO rooms (room) VALUES ($1);`
)

const (
	TableReadDepartment  = `SELECT id FROM departments WHERE code=$1;`
	TableReadDesignation = `SELECT id FROM faculty_designations WHERE designation=$1;`
	TableReadFaculty     = `SELECT id FROM faculty WHERE name=$1`
)
