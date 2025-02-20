<a id="readme-top"></a>

## What is <span style="color:#60a5fa">bluebeam</span>?

### AI-powered web research in one click, customized to your criteria.

bluebeam is a SaaS platform bridging market research and business intelligence.
It empowers you to conduct unlimited industry analyses by scanning websites of your choice against your own custom criteria, delivering tailored insights for smarter decision-making.

<title align="center">
  research websites in one click
  </br>
</title>

<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
      </ul>
    </li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#example">Example</a></li>
    <li><a href="#roadmap">Roadmap</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
    <li><a href="#contact">Contact</a></li>
  </ol>
</details>

### Built With

* [![Go][Go.dev]][Go-url]
* [![Docker][Docker.com]][Docker-url]
* [![Mariadb][Mariadb.org]][Mariadb-url]
* [![Javascript][Javascript.com]][Javascript-url]
* [![Python][Python.org]][Python-url]
* [![GoogleCloudPlatform][GoogleCloudPlatform.com]][GoogleCloudPlatform-url]
* OpenAI API and Gemini API

<!-- GETTING STARTED -->
## Getting Started

To get a local copy up and running, follow these simple example steps.
It will lauch the server of the application with docker-compose.
If you do not use a Linux distro, I recommand that you read the documentation on how to lauch docker-compose on your OS.

```sh
git clone https://github.com/BolvicBolvicovic/bluebeam
cd bluebeam/goserv
sudo make
```

Then, if you intend to use the extension, in a new terminal at the root of the repository, run the following commands:

```sh
cd ext
web-ext run
```

### Prerequisites

You will need go and docker-compose to lauch the project. If you intend to run the extension, you will need web-ext.
* go
```sh
wget https://go.dev/dl/go1.23.2.linux-amd64.tar.gz
rm -rf /usr/local/go && tar -C /usr/local -xzf go1.23.2.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin # You should put this in your .bashrc
```
* web-ext, docker and docker-compose
```sh
sudo apt update && sudo apt upgrade
sudo apt install -y web-ext docker docker-compose
```

If you want to use the google spreadsheet feature, add your google credentials file as googlecredentials.json at ./goserv/app/startup/
Furthermore, you will need to create an OAuth2.0 client in your google cloud account for your project and add the clientId to the dashboard.js file.

<p align="right">(<a href="#readme-top">back to top</a>)</p>


<!-- USAGE EXAMPLES -->
## Usage

### Once you have launched the <a href="#getting-started">server</a> 
- If you intend to use the web extension:<br/>
Running web-ext will open firefox. Because the server is running on localhost and the certificate is self-signed (at the moment),
you need to go to Settings -> Tools -> Advanced -> View Certificates -> Servers -> Add Exception then Add https://localhost.<br/>
- If you intend to use the dashboard:<br/>
Start your browser and go on localhost.<br/>
- If you intend to use the API:<br/>
You will have to wait, it is not implemented yet :)

### Register an account and log in

- If you intend to use the web extension:<br/>
Open your extension, (You have to be on a website that accepts scripting) and click on register a new account.
Fill all fields. It's important that the email you give is a valid Google email account. It will enable the extension to provide access to the output Google Spreadsheet.
If you get a response that is positive, you can go back and login with this account. Else, try with an other username/password.
- If you intend to use the dashboard:<br/>
Click on the login button at the top of the page and follow the same steps as for the extension. Once you are logged in, you will get access to a new section: dashboard.
<br/>
An HTTP-only cookie will keep you connected for a day and you will not need to reconnect yourself everytime you open the extension or the website.

### Analyse

If you have not given any criteria file, you will get an error. For both the extension and the website, this is done in the dashboard. (You have a button on the extension to access it directly.)
Current possible inputs are:
- a JSON file
- the first sheet of a google spreadsheet

For a JSON file, it has to contain an array of features described in the template <a href="https://github.com/BolvicBolvicovic/bluebeam/blob/main/example.json">example.json</a> at the root of the repository.
For the first sheet of a google spreadsheet, there is a <a href="images/example2.png">screenshot</a> in the folder images at the root of the repository.

Furthermore, you will need to put your OpenAI API key (mendatory for the extension) or your Gemini API key to make it work. It is stored with the data related to your account and only you has access to it.

Once everything is set, for the extension, you can click on the analyze button and you will get a response for the current page.
For the dashboard, you need to write the root urls you want to analyze then you click on the analyze button and you will get a response for the whole website (maximum size is currently 20 000 bytes but that can be manually change in the code).
<br/>
Current possible outputs are:
- a blob containing a JSON file
- a google spreadsheet

See roadmap below for future improvements.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Example

Here is a concret output example with the <a href="https://github.com/BolvicBolvicovic/bluebeam/blob/main/example.json">example.json</a> file used as an input and the <a href="https://go.dev/">Go website</a>.
![Alt text](example.png)

<!-- ROADMAP -->
## Roadmap

- [x] Response with a json/google spreadsheet that applies the chosen criteria on the website
- [x] Google spreadsheet format as the criteria's input
- [x] Possibility to audit many websites at the same time
- [x] Possibility to choose your AI
- [ ] API service that can be integrated to your application
- [ ] Scraping social-medias
- [ ] Response with the posibility to create a report with graphs and text

See the [open issues](https://github.com/BolvicBolvicovic/bluebeam/issues) for a full list of proposed features (and known issues).

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- CONTRIBUTING -->
## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request



<!-- LICENSE -->
## License

Distributed under the EUPL-1.2 License. See `LICENSE.txt` for more information.


<!-- CONTACT -->
## Contact

Project Link: [https://github.com/BolvicBolvicovic/bluebeam](https://github.com/BolvicBolvicovic/bluebeam)
Mail: victor.bolheme@gmail.com

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[license-shield]: https://img.shields.io/badge/license-EUPL%201.2-blue
[license-url]: https://github.com/BolvicBolvicovic/bluebeam/blob/main/LICENSE.txt
[linkedin-shield]: https://img.shields.io/badge/-LinkedIn-black.svg?style=for-the-badge&logo=linkedin&colorB=555
[linkedin-url]: https://linkedin.com/in/victorcornille
[product-screenshot]: images/screenshot.png
[Go.dev]: https://img.shields.io/badge/Go-00ADD8?logo=Go&logoColor=white&style=for-the-badge[Next-url]
[Go-url]: https://go.dev/
[Docker.com]: https://img.shields.io/badge/docker-257bd6?style=for-the-badge&logo=docker&logoColor=white
[Docker-url]: https://www.docker.com/
[Mariadb.org]: https://img.shields.io/badge/MariaDB-003545?style=for-the-badge&logo=mariadb&logoColor=white
[Mariadb-url]: https://mariadb.org/
[Javascript.com]: https://shields.io/badge/JavaScript-F7DF1E?logo=JavaScript&logoColor=000&style=flat-square
[Javascript-url]: https://www.javascript.com/
[Python.org]: https://img.shields.io/badge/python-3670A0?style=for-the-badge&logo=python&logoColor=ffdd54 
[Python-url]: https://www.python.org/
[GoogleCloudPlatform.com]: https://img.shields.io/badge/-Google%20Cloud%20Platform-4285F4?style=flat&logo=google%20cloud&logoColor=white 
[GoogleCloudPlatform-url]: https://console.cloud.google.com/ 
