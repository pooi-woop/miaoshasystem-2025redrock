package user

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"miaoshaSystem/global"
	"miaoshaSystem/sql"
	"net/http"
	"time"
)

type MyCustomClaims struct {
	Username           string `json:"username"` // 自定义字段，表示用户名
	jwt.StandardClaims        // 嵌入标准的JWT声明字段
}
type RequestData struct {
	Token   TokenData      `json:"token"`
	Product global.Product `json:"product"`
}

type TokenData struct {
	TokenString string `json:"tokenstring"`
}

func Register(c *gin.Context) {
	var user global.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := sql.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user creation failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user registered successfully"})
}

func Login(c *gin.Context) {
	var user struct {
		name string `json:"name"`
		pass string `json:"pass"`
	}
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var userr global.User
	sql.DB.Where("name = ? AND pass = ?", user.name, user.pass).First(&userr)
	var mySigningKey = []byte("mysecretkey")

	claims := MyCustomClaims{
		Username: userr.Name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
			Issuer:    "hym",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		fmt.Printf("Error signing token: %v\n", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

var user global.User

/*
	func Createmiaosha(c *gin.Context) {
		//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
		var tokenString string
		var requestData RequestData
		if err := c.ShouldBindJSON(&RequestData{}); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		tokenString = requestData.Token.TokenString
		var mySigningKey = []byte("mysecretkey")

		token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return mySigningKey, nil
		})

		if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {

			sql.DB.Where("name = ? ", claims.Username).First(&user)

		} else {
			fmt.Printf("Error validating token: %v\n", err)
			return
		}
		//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
		var Product global.Product

		Product = requestData.Product
		if err := sql.DB.Create(&Product).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Product creation failed"})
			return
		}
	}
*/
func Createmiaosha(c *gin.Context) {
	var requestData RequestData
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tokenString := requestData.Token.TokenString
	var mySigningKey = []byte("mysecretkey")

	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return mySigningKey, nil
	})

	if err != nil {
		fmt.Printf("Error parsing token: %v\n", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	if token == nil {
		fmt.Println("Token is nil")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		var user global.User
		sql.DB.Where("name = ?", claims.Username).First(&user)

		var product global.Product
		product = requestData.Product
		if err := sql.DB.Create(&product).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Product creation failed"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Product created successfully"})
	} else {
		fmt.Printf("Error validating token: %v\n", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}
}
func Miaosha(c *gin.Context) {
	productName := c.Param("productName")
	var tokenString string
	if err := c.ShouldBindJSON(&tokenString); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 鉴权逻辑
	var mySigningKey = []byte("mysecretkey")
	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return mySigningKey, nil
	})

	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		// 将秒杀请求发送到 Kafka 队列中
		err := global.SendToKafka(productName, claims.Username) //因为是抢购活动，一位用户不允许购买多台产品，所以不用传数量，
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request to Kafka"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Request accepted"})
	} else {
		fmt.Printf("Error validating token: %v\n", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}
}

/*接下来就该创建秒杀活动了，为了防止超卖，我得先减少库存数量再生成订单，这样即使订单未成功生成，也不会超卖
（当然，这样就得想办法防止少卖）*/
//为了处理高并发，我们使用kafka的消息队列
//我们用分布式🔒应对超卖问题
//因为只有一台电脑，所以不考虑负载均衡的问题（我也不会【逃ε=ε=ε┏(；ﾟロﾟ;)┛】）
//还有一个想法，我们可以用分布式架构保证秒杀这边炸了不会影响到登陆注册等其他功能
