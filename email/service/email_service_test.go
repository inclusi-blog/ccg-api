package service

import (
	"ccg-api/configuration"
	"ccg-api/constants"
	"ccg-api/email/email-client/email_client_request"
	mockEmailClient "ccg-api/email/email-client/mocks"
	"ccg-api/email/mocks"
	"ccg-api/email/models"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type emailServiceTestSuite struct {
	suite.Suite
	mockCtrl     *gomock.Controller
	context      *gin.Context
	recorder     *httptest.ResponseRecorder
	emailClient  *mockEmailClient.MockEmailClient
	emailConfig  *mocks.MockEmailClientConfig
	emailService EmailService
}

func TestEmailServiceTestSuite(t *testing.T) {
	suite.Run(t, new(emailServiceTestSuite))
}

func (suite *emailServiceTestSuite) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.recorder = httptest.NewRecorder()
	suite.context, _ = gin.CreateTestContext(suite.recorder)
	suite.context.Request, _ = http.NewRequest("GET", "some-url", nil)
	suite.emailClient = mockEmailClient.NewMockEmailClient(suite.mockCtrl)
	suite.emailConfig = mocks.NewMockEmailClientConfig(suite.mockCtrl)
	suite.emailService = NewEmailService(suite.emailClient, suite.emailConfig)
}

func (suite *emailServiceTestSuite) TearDownTest() {
	suite.mockCtrl.Finish()
}

func (suite emailServiceTestSuite) TestSendEmailShouldNotIncludeBaseTemplateAndSendEmailSuccessfully() {
	email := models.Email{
		From:    "gola@gola.xyz",
		To:      []string{"some@gmail.com"},
		Subject: "Hi!",
		Body: models.MessageBody{
			MimeType: "text/plain",
			Content:  "Hello User!",
		},
		Attachments: []models.Attachment{
			{
				FileName: "attachment1.pdf",
				Data:     []byte("Attachement1 Data!"),
			},
		},
	}

	expectedEmailClientRequest := email_client_request.EmailClientRequest{
		From:        email.From,
		To:          email.To,
		Subject:     email.Subject,
		Body:        email.Body,
		Attachments: email.Attachments,
	}

	suite.emailClient.EXPECT().Send(suite.context, &expectedEmailClientRequest)

	err := suite.emailService.Send(suite.context, email)
	suite.Nil(err)
}

func (suite emailServiceTestSuite) TestSendEmailShouldIncludeBaseTemplateAndSendEmailSuccessfully() {
	email := models.Email{
		From:    "gola@gola.xyz",
		To:      []string{"some@gmail.com"},
		Subject: "Hi!",
		Body: models.MessageBody{
			MimeType: "text/plain",
			Content:  "Hello User!",
		},
		Attachments: []models.Attachment{
			{
				FileName: "attachment1.pdf",
				Data:     []byte("Attachement1 Data!"),
			},
		},
		IncludeBaseTemplate: true,
	}

	suite.emailConfig.EXPECT().BaseTemplateFilePath().Return("../../email_templates/base_email_template.html")
	suite.emailConfig.EXPECT().LogoUrls().Return(configuration.LogoUrls{
		Mensuvadi:       "https://cdn.discordapp.com/attachments/731434048135757898/757125873030660096/unknown.png",
		Facebook:        "https://cdn.discordapp.com/attachments/731434048135757898/757125873030660096/unknown.png",
		Instagram:       "https://cdn.discordapp.com/attachments/731434048135757898/757125873030660096/unknown.png",
		Twitter:         "https://cdn.discordapp.com/attachments/731434048135757898/757125873030660096/unknown.png",
		LinkedIn:        "https://cdn.discordapp.com/attachments/731434048135757898/757125873030660096/unknown.png",
		DownloadIOS:     "https://cdn.discordapp.com/attachments/731434048135757898/757125873030660096/unknown.png",
		DownloadAndroid: "https://cdn.discordapp.com/attachments/731434048135757898/757125873030660096/unknown.png",
	})
	suite.emailConfig.EXPECT().OtherUrls().Return(configuration.Urls{
		HelpCenter:    "https://cdn.discordapp.com/attachments/731434048135757898/757125873030660096/unknown.png",
		PrivacyPolicy: "https://cdn.discordapp.com/attachments/731434048135757898/757125873030660096/unknown.png",
		Unsubscribe:   "https://cdn.discordapp.com/attachments/731434048135757898/757125873030660096/unknown.png",
		FAQUrl:        "https://cdn.discordapp.com/attachments/731434048135757898/757125873030660096/unknown.png",
	})
	suite.emailClient.EXPECT().Send(suite.context, gomock.Any()).Do(func(ctx *gin.Context, request *email_client_request.EmailClientRequest) {
		suite.Equal(email.From, request.From)
		suite.Equal(email.To, request.To)
		suite.Equal(email.Subject, request.Subject)
		suite.Equal(email.Body.MimeType, request.Body.MimeType)
		suite.Equal(email.Attachments, request.Attachments)
	}).Return(nil).Times(1)

	err := suite.emailService.Send(suite.context, email)
	suite.Nil(err)
}

func (suite emailServiceTestSuite) TestSendEmailShouldSendEmailReturnErrorIfClientUnableToSendEmail() {
	email := models.Email{
		From:    "gola@gola.xyz",
		To:      []string{"some@gmail.com"},
		Subject: "Hi!",
		Body: models.MessageBody{
			MimeType: "text/plain",
			Content:  "Hello User!",
		},
	}

	errMsg := "failed to send email"
	suite.emailClient.EXPECT().Send(suite.context, gomock.Any()).Return(errors.New(errMsg))
	expectedError := &constants.InternalServerError

	actualError := suite.emailService.Send(suite.context, email)
	suite.Equal(expectedError, actualError)
}

func (suite emailServiceTestSuite) TestSendEmailShouldIncludeBaseTemplateAndThrowErrorIfTemplateNotFound() {
	email := models.Email{
		From:    "gola@gola.xyz",
		To:      []string{"some@gmail.com"},
		Subject: "Hi!",
		Body: models.MessageBody{
			MimeType: "text/plain",
			Content:  "Hello User!",
		},
		Attachments: []models.Attachment{
			{
				FileName: "attachment1.pdf",
				Data:     []byte("Attachement1 Data!"),
			},
		},
		IncludeBaseTemplate: true,
	}

	suite.emailConfig.EXPECT().BaseTemplateFilePath().Return("base_email_template.html")

	err := suite.emailService.Send(suite.context, email)
	suite.Equal(&constants.InternalServerError, err)
}

func (suite emailServiceTestSuite) getBaseEmailTemplate() string {
	return `
    <html style="width:100%;font-family:arial, 'helvetica neue', helvetica, sans-serif;-webkit-text-size-adjust:100%;-ms-text-size-adjust:100%;padding:0;Margin:0;">
    <head>
        <meta charset="UTF-8">
        <meta content="width=device-width, initial-scale=1" name="viewport">
        <meta name="x-apple-disable-message-reformatting">
        <meta http-equiv="X-UA-Compatible" content="IE=edge">
        <meta content="telephone=no" name="format-detection">
        <!--[if (mso 16)]><style type="text/css">a {text-decoration: none;}</style><![endif]--><!--[if gte mso 9]><style>sup { font-size: 100% !important; }</style><![endif]--><!--[if gte mso 9]><style>sup { font-size: 100% !important; }</style><![endif]-->
        <style type="text/css">
            @media only screen and (max-width:600px) {p, ul li, ol li, a { font-size:14px!important; line-height:150%!important } h1 { font-size:30px!important; text-align:center; line-height:120%!important } h2 { font-size:26px!important; text-align:center; line-height:120%!important } h3 { font-size:20px!important; text-align:center; line-height:120%!important } h1 a { font-size:30px!important } h2 a { font-size:26px!important } h3 a { font-size:20px!important } .es-menu td a { font-size:16px!important } .es-header-body p, .es-header-body ul li, .es-header-body ol li, .es-header-body a { font-size:16px!important } .es-footer-body p, .es-footer-body ul li, .es-footer-body ol li, .es-footer-body a { font-size:16px!important } .es-infoblock p, .es-infoblock ul li, .es-infoblock ol li, .es-infoblock a { font-size:12px!important } *[class="gmail-fix"] { display:none!important } .es-m-txt-c, .es-m-txt-c h1, .es-m-txt-c h2, .es-m-txt-c h3 { text-align:center!important } .es-m-txt-r, .es-m-txt-r h1, .es-m-txt-r h2, .es-m-txt-r h3 { text-align:right!important } .es-m-txt-l, .es-m-txt-l h1, .es-m-txt-l h2, .es-m-txt-l h3 { text-align:left!important } .es-m-txt-r img, .es-m-txt-c img, .es-m-txt-l img { display:inline!important } .es-button-border { display:block!important } a.es-button { font-size:20px!important; display:block!important; border-left-width:0px!important; border-right-width:0px!important } .es-btn-fw { border-width:10px 0px!important; text-align:center!important } .es-adaptive table, .es-btn-fw, .es-btn-fw-brdr, .es-left, .es-right { width:100%!important } .es-content table, .es-header table, .es-footer table, .es-content, .es-footer, .es-header { width:100%!important; max-width:600px!important } .es-adapt-td { display:block!important; width:100%!important } .adapt-img { width:100%!important; height:auto!important } .es-m-p0 { padding:0px!important } .es-m-p0r { padding-right:0px!important } .es-m-p0l { padding-left:0px!important } .es-m-p0t { padding-top:0px!important } .es-m-p0b { padding-bottom:0!important } .es-m-p20b { padding-bottom:20px!important } .es-mobile-hidden, .es-hidden { display:none!important } .es-desk-hidden { display:table-row!important; width:auto!important; overflow:visible!important; float:none!important; max-height:inherit!important; line-height:inherit!important } .es-desk-menu-hidden { display:table-cell!important } table.es-table-not-adapt, .esd-block-html table { width:auto!important } table.es-social { display:inline-block!important } table.es-social td { display:inline-block!important } .es-m-margin { padding-left:5px!important; padding-right:5px!important; padding-top:5px!important; padding-bottom:5px!important } }
            #outlook a {
                padding:0;
            }
            .ExternalClass {
                width:100%;
            }
            .ExternalClass,
            .ExternalClass p,
            .ExternalClass span,
            .ExternalClass font,
            .ExternalClass td,
            .ExternalClass div {
                line-height:100%;
            }
            .es-button {
                mso-style-priority:100!important;
                text-decoration:none!important;
            }
            a[x-apple-data-detectors] {
                color:inherit!important;
                text-decoration:none!important;
                font-size:inherit!important;
                font-family:inherit!important;
                font-weight:inherit!important;
                line-height:inherit!important;
            }
            .es-desk-hidden {
                display:none;
                float:left;
                overflow:hidden;
                width:0;
                max-height:0;
                line-height:0;
                mso-hide:all;
            }
        </style>
    </head>
    <body style="width:100%;font-family:arial, 'helvetica neue', helvetica, sans-serif;-webkit-text-size-adjust:100%;-ms-text-size-adjust:100%;padding:0;Margin:0;">
    <div class="es-wrapper-color" style="background-color:#F6F6F6;">
        <!--[if gte mso 9]><v:background xmlns:v="urn:schemas-microsoft-com:vml" fill="t"><v:fill type="tile" color="#f6f6f6"></v:fill></v:background><![endif]-->
        <table class="es-wrapper" width="100%" cellspacing="0" cellpadding="0" style="mso-table-lspace:0pt;mso-table-rspace:0pt;border-collapse:collapse;border-spacing:0px;padding:0;Margin:0;width:100%;height:100%;background-repeat:repeat;background-position:center top;">
            <tr style="border-collapse:collapse;">
                <td class="es-m-margin" valign="top" style="padding:0;Margin:0;">
                    <table class="es-content" cellspacing="0" cellpadding="0" align="center" style="mso-table-lspace:0pt;mso-table-rspace:0pt;border-collapse:collapse;border-spacing:0px;table-layout:fixed !important;width:100%;">
                        <tr style="border-collapse:collapse;">
                            <td align="center" style="padding:0;Margin:0;">
                                <table class="es-content-body" width="900" cellspacing="0" cellpadding="0" bgcolor="#ffffff" align="center" style="mso-table-lspace:0pt;mso-table-rspace:0pt;border-collapse:collapse;border-spacing:0px;background-color:#FFFFFF;">
                                    <tr style="border-collapse:collapse;">
                                        <td align="left" style="padding:0;Margin:0;padding-top:20px;padding-left:20px;padding-right:20px;">
                                            <table cellpadding="0" cellspacing="0" width="100%" style="mso-table-lspace:0pt;mso-table-rspace:0pt;border-collapse:collapse;border-spacing:0px;">
                                                <tr style="border-collapse:collapse;">
                                                    <td width="860" align="center" valign="top" style="padding:0;Margin:0;">
                                                        <table cellpadding="0" cellspacing="0" width="100%" role="presentation" style="mso-table-lspace:0pt;mso-table-rspace:0pt;border-collapse:collapse;border-spacing:0px;">
                                                            <tr style="border-collapse:collapse;">
                                                                <td align="center" style="padding:0;Margin:0;padding-bottom:25px;font-size:0px;">
                                                                    <div style="display:flex;margin:1.25rem auto;justify-content:center;">
                                                                        <img src="http://10.176.6.53:7003/content/dam/optimus/BAS/email_images/logo_group.png" alt="logo"  width="120" height="51.5" style="display:block;border:0;outline:none;text-decoration:none;-ms-interpolation-mode:bicubic;width:7.5rem;height:3.44rem;">
                                                                    </div></td>
                                                            </tr>
                                                        </table></td>
                                                </tr>
                                            </table></td>
                                    </tr>
                                    <tr style="border-collapse:collapse;">
                                        <td align="left" style="padding:0;Margin:0;padding-left:20px;padding-right:20px;">
                                            <table cellpadding="0" cellspacing="0" width="100%" style="mso-table-lspace:0pt;mso-table-rspace:0pt;border-collapse:collapse;border-spacing:0px;">
                                                <tr style="border-collapse:collapse;">
                                                    <td width="860" align="center" valign="top" style="padding:0;Margin:0;">
                                                        Hello User!
                                                    </td>
                                                </tr>
                                            </table></td>
                                    </tr>
                                </table></td>
                        </tr>
                    </table>
                    <table cellpadding="0" cellspacing="0" class="es-content" align="center" style="mso-table-lspace:0pt;mso-table-rspace:0pt;border-collapse:collapse;border-spacing:0px;table-layout:fixed !important;width:100%;">
                        <tr style="border-collapse:collapse;">
                            <td align="center" style="padding:0;Margin:0;">
                                <table bgcolor="#ffffff" class="es-content-body" align="center" cellpadding="0" cellspacing="0" width="900" style="mso-table-lspace:0pt;mso-table-rspace:0pt;border-collapse:collapse;border-spacing:0px;background-color:#FFFFFF;">
                                    <tr style="border-collapse:collapse;">
                                        <td align="left" style="padding:0;Margin:0;padding-left:20px;padding-right:20px;">
                                            <table cellpadding="0" cellspacing="0" width="100%" style="mso-table-lspace:0pt;mso-table-rspace:0pt;border-collapse:collapse;border-spacing:0px;">
                                                <tr style="border-collapse:collapse;">
                                                    <td width="860" align="center" valign="top" style="padding:0;Margin:0;">
                                                        <table cellpadding="0" cellspacing="0" width="100%" role="presentation" style="mso-table-lspace:0pt;mso-table-rspace:0pt;border-collapse:collapse;border-spacing:0px;">
                                                            <tr style="border-collapse:collapse;">
                                                                <td align="center" style="padding:0;Margin:0;padding-bottom:10px;padding-top:20px;font-size:0;">
                                                                    <table border="0" width="100%" height="100%" cellpadding="0" cellspacing="0" role="presentation" style="mso-table-lspace:0pt;mso-table-rspace:0pt;border-collapse:collapse;border-spacing:0px;">
                                                                        <tr style="border-collapse:collapse;">
                                                                            <td style="padding:0;Margin:0px;border-bottom:2px solid #9D1D27;background:none;height:1px;width:100%;margin:0px;"></td>
                                                                        </tr>
                                                                    </table></td>
                                                            </tr>
                                                            <tr style="border-collapse:collapse;">
                                                                <td align="left" style="padding:0;Margin:0;">
                                                                    <div style="margin-top:10px;margin-bottom:10px;">
                                                                        <p style="Margin:0;-webkit-text-size-adjust:none;-ms-text-size-adjust:none;mso-line-height-rule:exactly;font-size:10px;font-family:arial, 'helvetica neue', helvetica, sans-serif;line-height:15px;color:#000000;">Disclaimer<br><br>Please note: IDFC FIRST Bank will never ask for any confidential information, do not share your username / password or Card details/CVV/OTP or PAN/AADHAAR details via email or via phone.<br><br>Beware of fraudulent emails that contain links of look-alike websites to mislead into entering sensitive information. IDFC FIRST Bank will never send such communication to customers asking for their personal or confidential information. Kindly visit IDFC FIRST Bank’s website instead of clicking on the links provided in emails from third parties.<br><br>These messages including any attachments are intended only for the addressee and may contain confidential, proprietary or legally privileged information. If you are not the named addressee or authorised to receive this mail, you shall not copy, forward, disclose or take any action based on this message or any part thereof. In such case, please notify the sender of this message and delete this message including any attachment of it from your computer system immediately. The recipient acknowledges that the views, opinions, conclusions and other information expressed in this message are those of the individual sender and shall be understood as neither given nor endorsed by IDFC FIRST Bank*, unless the sender does so expressly with due authority of IDFC FIRST Bank and IDFC FIRST Bank shall not be liable for any errors or omissions in the context of this message. Email transmission cannot be guaranteed to be secure or error-free as information could be intercepted, corrupted, lost, destroyed, arrive late or incomplete, or contain viruses. The sender therefore does not accept liability for any errors or omissions in the contents of this message, which arise as a result of e-mail transmission. This message, unless specifically stated in the email and followed by an agreement, does not tantamount to an offer or an acceptance of an offer by the sender.<br><br>*Includes IDFC FIRST Bank and all its subsidiary companies</p>
                                                                    </div></td>
                                                            </tr>
                                                            <tr style="border-collapse:collapse;">
                                                                <td align="center" style="padding:0;Margin:0;padding-top:10px;padding-bottom:10px;font-size:0;">
                                                                    <table border="0" width="100%" height="100%" cellpadding="0" cellspacing="0" role="presentation" style="mso-table-lspace:0pt;mso-table-rspace:0pt;border-collapse:collapse;border-spacing:0px;">
                                                                        <tr style="border-collapse:collapse;">
                                                                            <td style="padding:0;Margin:0px;border-bottom:2px solid #9D1D27;background:none;height:1px;width:100%;margin:0px;"></td>
                                                                        </tr>
                                                                    </table></td>
                                                            </tr>
                                                        </table></td>
                                                </tr>
                                            </table></td>
                                    </tr>
                                    <tr style="border-collapse:collapse;">
                                        <td align="left" style="Margin:0;padding-top:10px;padding-bottom:20px;padding-left:20px;padding-right:20px;">
                                            <!--[if mso]><table width="860" cellpadding="0" cellspacing="0" ><tr><td width="420" valign="top"><![endif]-->
                                            <table cellpadding="0" cellspacing="0" class="es-left" align="left" style="mso-table-lspace:0pt;mso-table-rspace:0pt;border-collapse:collapse;border-spacing:0px;float:left;">
                                                <tr style="border-collapse:collapse;">
                                                    <td width="420" class="es-m-p20b" align="left" style="padding:0;Margin:0;">
                                                        <table cellpadding="0" cellspacing="0" width="100%" role="presentation" style="mso-table-lspace:0pt;mso-table-rspace:0pt;border-collapse:collapse;border-spacing:0px;">
                                                            <tr style="border-collapse:collapse;">
                                                                <td align="left" style="padding:5px;Margin:0;">
                                                                    <div style="margin-top:10px;">
                                                                        <p style="Margin:0;-webkit-text-size-adjust:none;-ms-text-size-adjust:none;mso-line-height-rule:exactly;font-size:10px;font-family:arial, 'helvetica neue', helvetica, sans-serif;line-height:15px;color:#000000;">©2020 IDFC FIRST Bank. All rights reserved.</p>
                                                                    </div></td>
                                                            </tr>
                                                        </table></td>
                                                </tr>
                                            </table>
                                            <!--[if mso]></td><td width="20"></td><td width="420" valign="top"><![endif]-->
                                            <table cellpadding="0" cellspacing="0" class="es-right" align="right" style="mso-table-lspace:0pt;mso-table-rspace:0pt;border-collapse:collapse;border-spacing:0px;float:right;">
                                                <tr style="border-collapse:collapse;">
                                                    <td width="420" align="left" style="padding:0;Margin:0;">
                                                        <table cellpadding="0" cellspacing="0" width="100%" role="presentation" style="mso-table-lspace:0pt;mso-table-rspace:0pt;border-collapse:collapse;border-spacing:0px;">
                                                            <tr style="border-collapse:collapse;">
                                                                <td align="right" style="padding:0;Margin:0;font-size:0;">
                                                                    <div style="margin-top:10px;"></div>
                                                                    <table cellpadding="0" cellspacing="0" class="es-table-not-adapt es-social" role="presentation" style="mso-table-lspace:0pt;mso-table-rspace:0pt;border-collapse:collapse;border-spacing:0px;">
                                                                        <tr style="border-collapse:collapse;">
                                                                            <td align="center" valign="top" style="padding:0;Margin:0;padding-right:10px;"><a href="https://business.facebook.com/idfcfirstbank/" class="icon" style="-webkit-text-size-adjust:none;-ms-text-size-adjust:none;mso-line-height-rule:exactly;font-family:arial, 'helvetica neue', helvetica, sans-serif;font-size:14px;text-decoration:none;color:#9D1D27;"><img src="http://10.176.6.53:7003/content/dam/optimus/BAS/email_images/facebook.png" width="24" height="24" alt="facebook" style="width: 1.5rem;height: 1.5rem;margin-right: 5px;align-self: center;-ms-interpolation-mode:bicubic;"></a></td>
                                                                            <td align="center" valign="top" style="padding:0;Margin:0;padding-right:10px;"><a href="https://www.instagram.com/idfcfirstbank/" class="icon" style="-webkit-text-size-adjust:none;-ms-text-size-adjust:none;mso-line-height-rule:exactly;font-family:arial, 'helvetica neue', helvetica, sans-serif;font-size:14px;text-decoration:none;color:#9D1D27;"><img src="http://10.176.6.53:7003/content/dam/optimus/BAS/email_images/instagram.png" width="24" height="24" alt="instagram" style="width: 1.5rem;height: 1.5rem;margin-right: 5px;align-self: center;-ms-interpolation-mode:bicubic;"></a></td>
                                                                            <td align="center" valign="top" style="padding:0;Margin:0;padding-right:10px;"><a href="https://twitter.com/IDFCFIRSTBank" class="icon" style="-webkit-text-size-adjust:none;-ms-text-size-adjust:none;mso-line-height-rule:exactly;font-family:arial, 'helvetica neue', helvetica, sans-serif;font-size:14px;text-decoration:none;color:#9D1D27;"><img src="http://10.176.6.53:7003/content/dam/optimus/BAS/email_images/twitter.png" alt="twitter" width="24" height="24" style="width: 1.5rem;height: 1.5rem;margin-right: 5px;align-self: center;-ms-interpolation-mode:bicubic;"></a></td>
                                                                            <td align="center" valign="top" style="padding:0;Margin:0;padding-right:10px;"><a href="https://www.linkedin.com/company/idfcfirstbank" class="icon" style="-webkit-text-size-adjust:none;-ms-text-size-adjust:none;mso-line-height-rule:exactly;font-family:arial, 'helvetica neue', helvetica, sans-serif;font-size:14px;text-decoration:none;color:#9D1D27;"><img src="http://10.176.6.53:7003/content/dam/optimus/BAS/email_images/linked-in.png" alt="linkedin" width="24" height="24" style="width: 1.5rem;height: 1.5rem;margin-right: 5px;align-self: center;-ms-interpolation-mode:bicubic;"></a></td>
                                                                            <td align="center" valign="top" style="padding:0;Margin:0;"><a href="https://www.youtube.com/channel/UC3fyk0wieN6OdUIO-FARXDA" class="icon" style="-webkit-text-size-adjust:none;-ms-text-size-adjust:none;mso-line-height-rule:exactly;font-family:arial, 'helvetica neue', helvetica, sans-serif;font-size:14px;text-decoration:none;color:#9D1D27;"><img src="http://10.176.6.53:7003/content/dam/optimus/BAS/email_images/you-tube.png" alt="youtube" width="24" height="24" style="width: 1.5rem;height: 1.5rem;margin-right: 5px;align-self: center;-ms-interpolation-mode:bicubic;"></a></td>
                                                                        </tr>
                                                                    </table></td>
                                                            </tr>
                                                        </table></td>
                                                </tr>
                                            </table>
                                            <!--[if mso]></td></tr></table><![endif]-->
                                        </td>
                                    </tr>
                                    <tr style="border-collapse:collapse;">
                                        <td align="left" style="padding:0;Margin:0;padding-left:20px;padding-right:20px;">
                                            <!--[if mso]><table width="860" cellpadding="0" cellspacing="0"><tr><td width="314" valign="top"><![endif]-->
                                            <table cellpadding="0" cellspacing="0" class="es-left" align="left" style="mso-table-lspace:0pt;mso-table-rspace:0pt;border-collapse:collapse;border-spacing:0px;float:left;">
                                                <tr style="border-collapse:collapse;">
                                                    <td width="294" class="es-m-p0r es-m-p20b" align="center" style="padding:0;Margin:0;">
                                                        <table cellpadding="0" cellspacing="0" width="100%" role="presentation" style="mso-table-lspace:0pt;mso-table-rspace:0pt;border-collapse:collapse;border-spacing:0px;">
                                                            <tr style="border-collapse:collapse;">
                                                                <td width="21" style="padding:0;Margin:0;"><a href="tel:18004194332" style="-webkit-text-size-adjust:none;-ms-text-size-adjust:none;mso-line-height-rule:exactly;font-family:arial, 'helvetica neue', helvetica, sans-serif;font-size:14px;text-decoration:none;color:#9D1D27;">
                                                                        <div style="display:flex;font-size:10px;font-weight:bold;color:#000000;">
                                                                            <img src=http://10.176.6.53:7003/content/dam/optimus/BAS/email_images/call-icon.png alt="phone" width="16" height="24" style="width: 1rem;height: 1rem;align-self: center;-ms-interpolation-mode:bicubic;">
                                                                        </div></a>
                                                                </td>
                                                                <td>
                                                                    <a href="tel:18004194332" style="-webkit-text-size-adjust:none;-ms-text-size-adjust:none;mso-line-height-rule:exactly;font-family:arial, 'helvetica neue', helvetica, sans-serif;font-size:14px;text-decoration:none;color:#9D1D27;">
                                                                        <div style="display:flex;font-size:10px;font-weight:bold;color:#000000;">
                                                                            1800 419 4332
                                                                        </div>
                                                                    </a>
                                                                </td>
                                                            </tr>
                                                        </table></td>
                                                    <td class="es-hidden" width="20" style="padding:0;Margin:0;"></td>
                                                </tr>
                                            </table>
                                            <!--[if mso]></td><td width="263" valign="top"><![endif]-->
                                            <table cellpadding="0" cellspacing="0" class="es-left" align="left" style="mso-table-lspace:0pt;mso-table-rspace:0pt;border-collapse:collapse;border-spacing:0px;float:left;">
                                                <tr style="border-collapse:collapse;">
                                                    <td width="263" class="es-m-p20b" align="center" style="padding:0;Margin:0;">
                                                        <table cellpadding="0" cellspacing="0" width="100%" role="presentation" style="mso-table-lspace:0pt;mso-table-rspace:0pt;border-collapse:collapse;border-spacing:0px;">
                                                            <tr style="border-collapse:collapse;">
                                                                <td width="21" style="padding:0;Margin:0;"><a target="_blank" href="http://idfcfirstbank.com/" style="-webkit-text-size-adjust:none;-ms-text-size-adjust:none;mso-line-height-rule:exactly;font-family:arial, 'helvetica neue', helvetica, sans-serif;font-size:14px;text-decoration:none;color:#9D1D27;">
                                                                        <div style="display:flex;font-size:10px;font-weight:bold;color:#000000;">
                                                                            <img src=http://10.176.6.53:7003/content/dam/optimus/BAS/email_images/info.png alt="website" width="16" height="24" style="width: 1rem;height: 1rem;padding-right: 5px;align-self: center;-ms-interpolation-mode:bicubic;">
                                                                        </div></a>
                                                                </td>
                                                                <td>
                                                                    <a target="_blank" href="http://idfcfirstbank.com/" style="-webkit-text-size-adjust:none;-ms-text-size-adjust:none;mso-line-height-rule:exactly;font-family:arial, 'helvetica neue', helvetica, sans-serif;font-size:14px;text-decoration:none;color:#9D1D27;">
                                                                        <div style="display:flex;font-size:10px;font-weight:bold;color:#000000;">
                                                                            idfcfirstbank.com
                                                                        </div></a>
                                                                </td>
                                                            </tr>
                                                        </table></td>
                                                </tr>
                                            </table>
                                            <!--[if mso]></td><td width="20"></td><td width="263" valign="top"><![endif]-->
                                            <table cellpadding="0" cellspacing="0" class="es-right" align="right" style="mso-table-lspace:0pt;mso-table-rspace:0pt;border-collapse:collapse;border-spacing:0px;float:right;">
                                                <tr style="border-collapse:collapse;">
                                                    <td width="263" align="center" style="padding:0;Margin:0;">
                                                        <table cellpadding="0" cellspacing="0" width="100%" role="presentation" style="mso-table-lspace:0pt;mso-table-rspace:0pt;border-collapse:collapse;border-spacing:0px;">
                                                            <tr style="border-collapse:collapse;">
                                                                <td width="29" style="padding:0;Margin:0;"><a href="mailto:banker@idfcfirstbank.com" style="-webkit-text-size-adjust:none;-ms-text-size-adjust:none;mso-line-height-rule:exactly;font-family:arial, 'helvetica neue', helvetica, sans-serif;font-size:14px;text-decoration:none;color:#9D1D27;">
                                                                        <div style="display:flex;font-size:10px;font-weight:bold;color:#000000;">
                                                                            <img src=http://10.176.6.53:7003/content/dam/optimus/BAS/email_images/email.png alt="email" width="24" height="24" style="width: 1.5rem;height: 1rem;align-self: center;-ms-interpolation-mode:bicubic;">
                                                                        </div></a>
                                                                </td>
                                                                <td>
                                                                    <a href="mailto:banker@idfcfirstbank.com" style="-webkit-text-size-adjust:none;-ms-text-size-adjust:none;mso-line-height-rule:exactly;font-family:arial, 'helvetica neue', helvetica, sans-serif;font-size:14px;text-decoration:none;color:#9D1D27;">
                                                                        <div style="display:flex;font-size:10px;font-weight:bold;color:#000000;">
                                                                            banker@idfcfirstbank.com
                                                                        </div></a>
                                                                </td>
                                                            </tr>
                                                        </table></td>
                                                </tr>
                                            </table>
                                            <!--[if mso]></td></tr></table><![endif]-->
                                        </td>
                                    </tr>
                                </table></td>
                        </tr>
                    </table></td>
            </tr>
        </table>
    </div>
    </body>
    </html>
`
}
