package com.example;

import org.openqa.selenium.By;
import org.openqa.selenium.OutputType;
import org.openqa.selenium.TakesScreenshot;
import org.openqa.selenium.WebDriver;
import org.openqa.selenium.WebElement;
import org.openqa.selenium.chrome.ChromeOptions;
import org.openqa.selenium.remote.RemoteWebDriver;

import java.io.File;
import java.net.URI;
import java.net.URL;
import java.nio.file.Files;

public class Test {

    public static By byText(String text) {
        return By.xpath ("//*[contains(text(),'"+text+"')]");
    }

    public static void main(String[] args) {
        try {
            // Connect to the Selenium server running in Docker
            URL url = new URI("http://selenium:4444/wd/hub").toURL();
//            DesiredCapabilities capabilities = DesiredCapabilities.();
            WebDriver driver = new RemoteWebDriver(url, new ChromeOptions());

            // Navigate to Google
            driver.get("http://maincontainer:8080");

            // Find the search box and perform a search
            WebElement bt = driver.findElement(byText("Button 2"));
            if (!driver.findElements(byText("Button 3")).isEmpty()) {
                throw new Exception("Button 3 already exists");
            }

            bt.click();
            WebElement bt1 = driver.findElement(byText("Button 3"));
            bt1.click();
            System.out.println("Passed!");
            File scrFile = ((TakesScreenshot)driver).getScreenshotAs(OutputType.FILE);
            File f = new File("screenshots/picture1.png");
            Files.move(scrFile.toPath(), f.toPath());
            System.out.println("Wrote file: " + f.getAbsolutePath());
            // Close the browser
            driver.quit();
        } catch (Exception e) {
            e.printStackTrace();
            System.exit(1); // error state.
        }
    }
}
