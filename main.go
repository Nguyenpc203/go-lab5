package main

import (
	"net/http"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

// Định nghĩa cấu trúc dữ liệu cho Student.
type Student struct {
	gorm.Model        // Tự động thêm các trường như ID, CreatedAt, UpdatedAt, DeletedAt
	Name       string `json:"name"`
	Age        int    `json:"age"`
	Class      string `json:"class"`
	Email      string `json:"email"`
}

var db *gorm.DB
var err error

func main() {
	// Kết nối tới cơ sở dữ liệu MySQL.
	dsn := "root:passroot@tcp(127.0.0.1:3308)/student_management?parseTime=true&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Tự động tạo bảng students nếu chưa tồn tại.
	db.AutoMigrate(&Student{})

	// Tạo router với Gin.
	router := gin.Default()

	// Các endpoint thực hiện CRUD.
	router.POST("/students", createStudent)       // Thêm mới sinh viên
	router.GET("/students", getStudents)          // Lấy danh sách sinh viên
	router.GET("/students/:id", getStudentByID)   // Lấy thông tin sinh viên theo ID
	router.PUT("/students/:id", updateStudent)    // Cập nhật thông tin sinh viên theo ID
	router.DELETE("/students/:id", deleteStudent) // Xóa sinh viên theo ID

	// Khởi chạy server trên port 8080.
	router.Run(":8080")
}

// Hàm thêm mới sinh viên.
func createStudent(c *gin.Context) {
	var student Student
	if err := c.ShouldBindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.Create(&student)
	c.JSON(http.StatusOK, student)
}

// Hàm lấy tất cả sinh viên.
func getStudents(c *gin.Context) {
	var students []Student
	db.Find(&students)
	c.JSON(http.StatusOK, students)
}

// Hàm lấy sinh viên theo ID.
func getStudentByID(c *gin.Context) {
	id := c.Param("id")
	var student Student
	result := db.First(&student, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}
	c.JSON(http.StatusOK, student)
}

// Hàm cập nhật thông tin sinh viên theo ID.
func updateStudent(c *gin.Context) {
	id := c.Param("id")
	var student Student
	if err := db.First(&student, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	// Ràng buộc dữ liệu từ yêu cầu JSON tới student.
	if err := c.ShouldBindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.Save(&student)
	c.JSON(http.StatusOK, student)
}

// Hàm xóa sinh viên theo ID.
func deleteStudent(c *gin.Context) {
	id := c.Param("id")
	var student Student
	if err := db.First(&student, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}
	db.Delete(&student)
	c.JSON(http.StatusOK, gin.H{"message": "Student deleted successfully"})
}
