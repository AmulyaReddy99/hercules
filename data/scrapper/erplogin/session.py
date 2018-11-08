import requests

ERP_HOME_URL = 'https://erp.iitkgp.ac.in/IIT_ERP3/'
SECURITY_QUESTION_URL = 'https://erp.iitkgp.ac.in/SSOAdministration/getSecurityQues.htm'
LOGIN_URL = 'https://erp.iitkgp.ac.in/SSOAdministration/auth.htm'
GET_ACAD_TOKEN_URL = 'https://erp.iitkgp.ac.in/Acad/central_breadth_tt.jsp?action=second'

class ERPSession:
    """
    An erp session

    Logs into erp after prompting the user to enter details

    Class Attributes::
    -self.sessionToken #the erp SessionToken or JSESSIONID (as named by erp)
    -self.ssoToken    #the SSOToken obtained after loggin into erp
    -self.academicToken #the academic token, useful for getting the data from erp

    Basic Usage::

    >>>from erplogin.session import ERPSession
    >>>s = ERPSession(roll_no, password)
    >>>question = s.get_security_question() # Returns the security question from ERP
    # User feeds the answer to security question.
    >>>s.LoginERP(answer)

    Obtain the cookies by accesing the corresponding attributes.
    """

    headers = {
        'timeout': '20',
        'User-Agent': 'Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Ubuntu Chromium/51.0.2704.79 Chrome/51.0.2704.79 Safari/537.36',
    }

    def __init__(self, roll_no, password):
        self.sess = requests.Session()
        self.roll_no = roll_no
        self._password = password
        self.sessionToken = ERPSession.__generate_session_cookie(self)
        self.question_answer = ''
        self.academicToken = ''
        self.ssoToken = ''
    
    def __generate_session_cookie (self):
        """
        Automatically called upon declaration
        Requests for the session id by requesting for erp home page.

        Arguments:
        > NIL

        Returns 
        > sessionToken
        """

        response_erp = self.sess.get(ERP_HOME_URL)
        if response_erp.status_code != 200:
            return None

        return (response_erp.cookies['JSESSIONID'])


    def get_security_question(self):
        """
        Request the security question from ERP. 

        Argument: 
        > NIL

        Returns :
        > question # Security question from ERP
        """
        response_security_question = self.sess.post (SECURITY_QUESTION_URL, data = {'user_id': self.roll_no},headers=self.headers)

        if(response_security_question.status_code == 200):
            if(response_security_question.text != 'FALSE'):
                return response_security_question.text

        return None

    def LoginERP(self,answer):
        """
        Logs into ERP and sets the academicToken and SSOToken accordingly

        Argument:
        > answer # answer to the security question

        >Returns the 
        """
        login_details = {
                        'user_id': self.roll_no,
                        'password': self._password,
                        'answer': answer,
                        'sessionToken': self.sessionToken,
                        'requestedUrl': ERP_HOME_URL,
                    }

        response_login_auth_htm = self.sess.post(LOGIN_URL
                                            , data=login_details
                                            , headers=self.headers
                                            )

        if (response_login_auth_htm.status_code == 200 \
            and len(response_login_auth_htm.history) == 2):

            ssoToken = response_login_auth_htm.history[1].cookies['ssoToken']
            self.ssoToken = ssoToken

            response_acad = self.sess.get(GET_ACAD_TOKEN_URL)
            if(response_acad.status_code == 200):
                try:
                    self.academicToken = response_acad.cookies["JSID#/Acad"]
                except:
                    self.academicToken = None
                return

            self.academicToken = None
            return

        self.ssoToken = None
        self.academicToken = None
