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
    public static void main(String[] args) {
        try {
            // Connect to the Selenium server running in Docker
            URL url = new URI("http://selenium:4444/wd/hub").toURL();
//            DesiredCapabilities capabilities = DesiredCapabilities.();
            WebDriver driver = new RemoteWebDriver(url, new ChromeOptions());

            // Navigate to Google
            driver.get("http://maincontainer:8080");

            // Find the search box and perform a search
            WebElement searchBox = driver.findElement(By.name("q"));
            searchBox.sendKeys("Selenium Docker");
            searchBox.submit();

            // Wait for the results to load
            Thread.sleep(2000);

            // Print the title of the page
            System.out.println("Title: " + driver.getTitle());
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
