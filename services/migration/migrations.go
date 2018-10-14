package migration

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"

	"github.com/kshitij10496/hercules/common"
)

func containsDesignations(designations []common.FacultyDesignation, d common.FacultyDesignation) bool {
	for _, designation := range designations {
		if d == designation {
			return true
		}
	}
	return false
}

func containsDepartments(departments []readDepartment, d readDepartment) bool {
	for _, department := range departments {
		if d == department {
			return true
		}
	}
	return false
}

func createSetDepartments(departments []readDepartment) []readDepartment {
	setDepartments := []readDepartment{}
	for _, d := range departments {
		if !containsDepartments(setDepartments, d) {
			setDepartments = append(setDepartments, d)
		}
	}
	return setDepartments
}

func createSetDesignations(designations []common.FacultyDesignation) []common.FacultyDesignation {
	setDesignations := []common.FacultyDesignation{}
	for _, d := range designations {
		if !containsDesignations(setDesignations, d) {
			setDesignations = append(setDesignations, d)
		}
	}
	return setDesignations
}

type readDepartment struct {
	Name string `json:"department"`
	Code string `json:"code"`
}

func readFromJSONDepartments(db *sql.DB, filename string) error {
	// Open JSON file
	departmentsFile, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer departmentsFile.Close()

	// Load data into a "Departments" value
	decoder := json.NewDecoder(departmentsFile)
	var departments []readDepartment
	err = decoder.Decode(&departments)
	if err != nil {
		return err
	}

	// Insert data into db
	// TODO: Implement a SQL transaction here.
	departmentsSet := createSetDepartments(departments)
	log.Println("Department Set:", departmentsSet)
	for _, department := range departmentsSet {
		_, err := db.Exec(common.TableInsertionDepartment, department.Code, department.Name)
		if err != nil {
			log.Println("[insertion] departments:", department, err)
		}
	}

	return nil
}

func readFromJSONFacultyDesignations(db *sql.DB, filename string) error {
	designationsFile, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer designationsFile.Close()

	decoder := json.NewDecoder(designationsFile)
	var designations []common.FacultyDesignation
	err = decoder.Decode(&designations)
	if err != nil {
		return err
	}

	designationsSet := createSetDesignations(designations)
	log.Println("DESIGNATION SET:", designationsSet)
	for _, designation := range designationsSet {
		_, err := db.Exec(common.TableInsertionDesignation, designation.Designation)
		if err != nil {
			log.Println("[insertion] faculty_designations:", designation, err)
		}
	}
	return nil
}

type readFacultyMember struct {
	Name        string `json:"faculty"`
	Department  string `json:"department"`
	Designation string `json:"designation"`
	DeptCode    string `json:"code"`
}

type readFaculty []readFacultyMember

func readFromJSONFaculty(db *sql.DB, filename string) error {
	facultyFile, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer facultyFile.Close()

	decoder := json.NewDecoder(facultyFile)
	var faculty readFaculty
	err = decoder.Decode(&faculty)
	if err != nil {
		return err
	}

	for _, member := range faculty {
		var designationID, departmentID int

		row := db.QueryRow(common.TableReadDesignation, member.Designation)
		if err := row.Scan(&designationID); err != nil {
			log.Println("[read] faculty_desingations:", member.Designation, err)
			continue
		}

		row = db.QueryRow(common.TableReadDepartment, member.DeptCode)
		if err := row.Scan(&departmentID); err != nil {
			log.Println("[read] departments:", member.DeptCode, err)
			continue
		}

		_, err = db.Exec(common.TableInsertionFaculty, member.Name, designationID, departmentID)
		if err != nil {
			log.Println("[insertion] faculty:", member, err)
			continue
		}
	}
	return nil
}
