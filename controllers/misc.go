package controllers

import (
	. "donutBackend/logger"
	emailsender "donutBackend/models/emailSender"
	"donutBackend/models/messages"
	organization "donutBackend/models/orgs"
	pendingEmail "donutBackend/models/pendingEmails"
	"donutBackend/models/users"
	"donutBackend/models/waitlist"
	weblinks "donutBackend/models/web_links"
	"donutBackend/utils/mail"
	email "donutBackend/utils/mail"
	"fmt"
	"strings"

	//"encoding/base64"
	"encoding/json"
	"net/http"

	//"gopkg.in/robfig/cron.v2"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func JoinWaitlist(c *gin.Context) {
	var waitlistedUser waitlist.WaitlistedUser
	if err := c.ShouldBindBodyWith(&waitlistedUser, binding.JSON); err != nil {
		Logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to Add to Waitlist",
			"error":   err.Error(),
		})
		return
	}

	firstName := strings.Split(waitlistedUser.Name, " ")[0]
	//count no of chars in firstName
	if len(firstName) > 5 {
		firstName = strings.ToUpper(firstName[:5])
	} else {
		firstName = strings.ToUpper(firstName)
	}

	if _, err := waitlist.Insert(waitlistedUser); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"message": "Failed To Add to Waitlist",
			"error":   "Already added to the Waitlist",
		})
		return
	}

	link, err := weblinks.AddLink(waitlistedUser.Email, "FALSE")
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"message": "Failed To Add to Waitlist",
			"error":   "Already added to the Waitlist AddLink",
		})
		return
	}

	linkId := link.Id[:5]
	//linkId:="ABCD"
	referral := firstName + "-" + linkId

	count, err := waitlist.GetCount()
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"message": "Failed To Add to Waitlist",
			"error":   "Already added to the Waitlist Count",
		})
		return
	}

	count = count + 1020
	//count:=1030
	//convert count to string
	countString := fmt.Sprintf("%d", count)

	waitlistEmail := `<div class="background center" style="background-color: #FFF8FC;padding: 36px 32px 64px;border-radius: 20px;max-width: 400px;margin-bottom: 20px;margin-left: auto;margin-right: auto;">
	<img class="with-space" width="40" height="40" src="https://imagedelivery.net/_ytNarHSFtZWvy0qTLubhg/3f5bf187-93df-4a98-8c7d-567240719f00/public" alt="Donut Logo" style="margin-bottom: 36px;"/>
  
	<div class="poppins-semibold heading-text no-space" style="margin-bottom: 4px;font-family: 'Poppins', 'Arial', 'sans-serif';font-weight: 600;font-size: 16px;">Welcome to awesomeness!</div>
	<div class="poppins-regular body-text with-space" style="margin-bottom: 36px;font-family: 'Poppins', 'Arial', 'sans-serif';font-weight: 400;font-size: 14px;">Thanks for joining the Donut waitlist. You’re in for so much fun! 🥳</div>
  
	<div class="with-space" style="margin-bottom: 36px;">
	  <img height="100%" width="100%" src="https://imagedelivery.net/_ytNarHSFtZWvy0qTLubhg/92e134e0-7975-4ee6-c92e-228f8d9a7d00/public" alt="We are building an ecosystem to make your organisation 10X better" style="object-fit: contain;">
	</div>
	<!--    <div class="container better with-space">-->
	<!--        <div class="poppins-regular body-text text-center">We are building an ecosystem to make your organisation</div>-->
	<!--        <div class="poppins-semibold display-text text-center text-pink ten-x">10<span class="x">x</span></div>-->
	<!--        <div class="poppins-regular body-text text-center">better</div>-->
	<!--    </div>-->
  
	<div class="with-space" style="margin-bottom: 36px;">
	  <table class="no-border" border="0" cellpadding="0" cellspacing="0px" style="border-spacing: 0;border-collapse: collapse;">
		<tr>
		  <td>
			<div class="poppins-semibold body-text no-space" style="margin-bottom: 4px;font-family: 'Poppins', 'Arial', 'sans-serif';font-weight: 600;font-size: 14px;">
			  Enhance your<br>online presence
			</div>
			<div class="poppins-regular body-text" style="font-family: 'Poppins', 'Arial', 'sans-serif';font-weight: 400;font-size: 14px;">
			  Donut helps you create a tailor-made public profile that you can share with people
			</div>
		  </td>
		  <td>
			<img style="margin-left: 10px;object-fit: contain;" height="100%" width="100%" src="https://imagedelivery.net/_ytNarHSFtZWvy0qTLubhg/39bbf2cd-c07a-41fa-4652-57efa3160000/public" alt="Profile">
		  </td>
		</tr>
	  </table>
	</div>
  
	<div class="with-space" style="margin-bottom: 36px;"><img height="100%" width="100%" src="https://imagedelivery.net/_ytNarHSFtZWvy0qTLubhg/07d14a17-4ce6-4a6c-44bd-19b86718b100/public" alt="Analytics" style="object-fit: contain;"></div>
  
	<div class="poppins-semibold body-text text-right no-space" style="margin-bottom: 4px;text-align: right;font-family: 'Poppins', 'Arial', 'sans-serif';font-weight: 600;font-size: 14px;">Campaigns, Tools, Analytics, and more!</div>
	<div class="poppins-regular body-text text-right with-space" style="margin-bottom: 36px;text-align: right;font-family: 'Poppins', 'Arial', 'sans-serif';font-weight: 400;font-size: 14px;">In your profile, you will find tools to help your organisation raise funds, get volunteers and run your organisation at scale.</div>
  
	<div class="container waitlist with-space" style="margin-bottom: 36px;background-color: #FFFFFF;border-radius: 28px;padding: 24px 32px 55px;">
	  <div class="poppins-regular body-text text-center text-grey with-little-space" style="margin-bottom: 24px;text-align: center;color: #8E8E8E;font-family: 'Poppins', 'Arial', 'sans-serif';font-weight: 400;font-size: 14px;">Your current position on the waitlist is</div>
	  <div class="poppins-semibold display-text text-center waitlist-number with-little-space" style="margin-bottom: 24px;text-align: center;font-family: 'Poppins', 'Arial', 'sans-serif';font-weight: 600;font-size: 36px;color: #FE7088;">#` + countString + `</div>
	  <div class="poppins-semibold body-text text-center no-space" style="margin-bottom: 4px;text-align: center;font-family: 'Poppins', 'Arial', 'sans-serif';font-weight: 600;font-size: 14px;">Want exclusive early access?</div>
	  <div class="poppins-regular sub-text text-center with-little-space" style="margin-bottom: 24px;text-align: center;font-family: 'Poppins', 'Arial', 'sans-serif';font-weight: 400;font-size: 13px;">Share the awesomeness with your friends & network and cut the line short.</div>
	  <a href="https://donut.gives/refer?tag=` + link.Id + `&refer=` + referral + `" target="_blank" style="text-decoration: none">
		<div class="primary-button" style="background-color: #FF6A85;border-radius: 20px;height: 40px;color: white;align-items: center;align-content: center;text-align: center;line-height: 38px;font-family: 'Clash Grotesk', 'Space Grotesk', 'Arial', 'sans-serif';letter-spacing: 1px;font-size: 14px;">
		  <table class="center no-border" border="0" cellpadding="0" cellspacing="0" style="margin-left: auto;margin-right: auto;border-spacing: 0;border-collapse: collapse;">
			<tr>
			  <td>Share referral link</td>
			  <td><img class="arrow" height="20px" src="https://imagedelivery.net/_ytNarHSFtZWvy0qTLubhg/a98fc695-62d3-4fe5-c895-24d52c24c900/public" alt="Open" style="object-fit: contain;margin-top: 9px;margin-bottom: -3px;margin-left: 10px;"></td>
			</tr>
		  </table>
		</div>
	  </a>
	</div>
  
	<div class="poppins-regular body-text text-center with-space" style="margin-bottom: 36px;text-align: center;font-family: 'Poppins', 'Arial', 'sans-serif';font-weight: 400;font-size: 14px;">
	  We will contact you shortly for easy onboarding soon. If you have any queries, we are just a mail away. Feel free to reply to this mail to get an assured reply from the team.
	</div>
  
	<div class="container story" style="background-color: #FFFFFF;border-radius: 28px;padding: 28px 32px 36px;">
	  <div class="poppins-semibold body-text text-center with-very-little-space" style="margin-bottom: 12px;text-align: center;font-family: 'Poppins', 'Arial', 'sans-serif';font-weight: 600;font-size: 14px;">Want to know our story?</div>
	  <a href="https://donutgives.notion.site/Donut-fc5fcf1a735e41a38d0332cdbcf787bb" target="_blank" style="text-decoration: none"><div class="secondary-button" style="border-radius: 22px;height: 40px;border: 2px solid #FF6A85;color: #FF6A85;text-align: center;line-height: 38px;font-family: 'Clash Grotesk', 'Space Grotesk', 'Arial', 'sans-serif';letter-spacing: 1px;font-size: 14px;">Read here</div></a>
	</div>
  </div>
  `

	go func() {
		//send email atmost 3 times and unitl sent
		sent := false
		for i := 0; i < 1; i++ {

			subject := "Welcome to the Waitlist! 🎉 Invite Your Friends & Get Early Access"

			//subject := "Welcome to the waitlist! "

			err := email.SendMail(waitlistedUser.Email, subject, "text/html", waitlistEmail)
			if err == nil {
				sent = true
				break
			}
			Logger.Errorf("1Failed to send email to %s: %v", waitlistedUser.Email, err)

			doBreak := false
			err = email.RefreshAccessToken()
			if err == nil {
				err := email.SendMail(waitlistedUser.Email, subject, "text/html", waitlistEmail)
				if err == nil {
					sent = true
					doBreak = true
					break
				}
				Logger.Errorf("2Failed to send email to %s: %v", waitlistedUser.Email, err)
			} else {
				Logger.Errorf("Error refreshing token: %v", err)
			}
			if doBreak {
				break
			}
			Logger.Errorf("Sent= %s doBreak= %s", sent, doBreak)

			Logger.Errorf("Failed to send email to %s: %v", waitlistedUser.Email, err)
			Logger.Errorf("Failed to refresh access token for email %s: %v", mail.Email, err)
			emailsender.SetDeactivated(mail.Email)
			err = mail.SendMailBySMTP("dev.donut.gives@gmail.com", "Current Email Sender Deactivated", "text/plain", "Please login for gmail credentials again.")
			if err != nil {
				Logger.Errorf("Failed to send email for deactivated email sender %s: %v", mail.Email, err)
			}

		}
		if !sent {

			pending := pendingEmail.PendingEmail{
				Email:  waitlistedUser.Email,
				Reason: "Waitlist Email",
				Body:   waitlistEmail,
			}
			_, err := pendingEmail.Insert(pending)
			if err != nil {
				Logger.Errorf("Failed to insert pending email %s while %s: %v", waitlistedUser.Email, "waitlisting", err)
			}
		}
	}()

	// err = email.SendMail(waitlistedUser.Email,"Congratulation On Being Waitlisted!","text/html",waitlistEmail)
	// if err != nil {
	// 	c.JSON(http.StatusBadGateway, gin.H{
	// 		"message": "Failed to Send Waitlist Email",
	// 		"error":   err.Error(),
	// 	})
	// 	return
	// }

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully Added to Waitlist",
	})
	return

}

func ContactUs(c *gin.Context) {
	var message messages.Message
	if err := c.ShouldBindBodyWith(&message, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to Add to Save Message",
			"error":   err.Error(),
		})
	}

	_, err := messages.Insert(message)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"message": "Failed to Save Message",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully Sent the Message",
	})

}

func DiscordInvite(c *gin.Context) {
	c.Redirect(http.StatusTemporaryRedirect, "https://discord.gg/gXPA9xeFw8")
}

func GetProfile(c *gin.Context) {

	params := c.Request.URL.Query()
	email := params.Get("email")
	target := params.Get("target")

	if email == "" {

		entity := c.GetString("request")

		if entity == "user" {

			user := users.GoogleUser{}
			err := json.Unmarshal([]byte(c.GetString("user")), &user)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}

			userProfile, err := users.GetUserProfile(user.Email)
			if err != nil {
				c.JSON(http.StatusBadGateway, gin.H{
					"message": "Failed to Get User Profile",
					"error":   err.Error(),
				})
				return
			}

			returnJSON := struct {
				Profile  users.GoogleUserProfile
				Verified string
			}{
				Profile:  userProfile,
				Verified: "true",
			}

			c.JSON(http.StatusOK, gin.H{
				"message": "Successfully Got User Profile",
				"data":    returnJSON,
			})

		} else if entity == "org" {

			org := organization.Organization{}
			err := json.Unmarshal([]byte(c.GetString("org")), &org)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": err.Error(),
				})
				return
			}

			orgProfile, err := organization.GetOrgProfile(org.Email)
			if err != nil {
				c.JSON(http.StatusBadGateway, gin.H{
					"message": "Failed to Get Org Profile",
					"error":   err.Error(),
				})
				return
			}

			returnJSON := struct {
				Profile  organization.OrganizationProfile
				Verified string
			}{
				Profile:  orgProfile,
				Verified: "true",
			}

			c.JSON(http.StatusOK, gin.H{
				"message": "Successfully Got Org Profile",
				"data":    returnJSON,
			})

		}

	} else {

		if target == "user" {

			userProfile, err := users.GetUserProfile(email)
			if err != nil {
				c.JSON(http.StatusBadGateway, gin.H{
					"message": "Failed to Get User Profile",
					"error":   err.Error(),
				})
				return
			}

			returnJSON := struct {
				Profile  users.GoogleUserProfile
				Verified string
			}{
				Profile:  userProfile,
				Verified: "false",
			}

			c.JSON(http.StatusOK, gin.H{
				"message": "Successfully Got User Profile",
				"data":    returnJSON,
			})

		} else if target == "org" {

			orgProfile, err := organization.GetOrgProfile(email)
			if err != nil {
				c.JSON(http.StatusBadGateway, gin.H{
					"message": "Failed to Get Org Profile",
					"error":   err.Error(),
				})
				return
			}

			returnJSON := struct {
				Profile  organization.OrganizationProfile
				Verified string
			}{
				Profile:  orgProfile,
				Verified: "false",
			}

			c.JSON(http.StatusOK, gin.H{
				"message": "Successfully Got Org Profile",
				"data":    returnJSON,
			})
		}
	}

}
