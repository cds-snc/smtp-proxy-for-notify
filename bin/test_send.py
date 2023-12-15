# Send an email to a local smtp server using the smtplib module
import os
import random
import smtplib
import ssl
import string
import email.utils
from email.mime.text import MIMEText
from email.message import EmailMessage

# Generate random string for the attachment
def randomword(length):
    letters = string.ascii_lowercase
    return ''.join(random.choice(letters) for i in range(length))

email = EmailMessage()
email['subject'] = "Simple test message"
email.set_content("This is the body of the message.")

attachment_size_in_kb = 1000

#for i in range(5):
#    email.add_attachment(
#      randomword(attachment_size_in_kb * 1024).encode('utf-8'),
#        filename="attachment-{}.txt".format(i),
#        maintype="application",
#        subtype="txt"
#    )

context = ssl.SSLContext(ssl.PROTOCOL_TLSv1_2)
server = smtplib.SMTP_SSL("0.0.0.0", 1025)
server.ehlo()  # send the extended hello to our server
server.set_debuglevel(True) # show communication with the server
try:
    server.login(os.environ["SMTP_USERNAME"], os.environ["SMTP_PASSWORD"])
    server.sendmail(os.environ["TEST_SENDER"], [os.environ["TEST_RECIPIENT"]], email.as_string())
finally:
    server.quit()
