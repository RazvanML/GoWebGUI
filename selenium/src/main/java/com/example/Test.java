package com.example;

import org.openqa.selenium.By;
import org.openqa.selenium.OutputType;
import org.openqa.selenium.TakesScreenshot;
import org.openqa.selenium.WebDriver;
import org.openqa.selenium.WebElement;
import org.openqa.selenium.chrome.ChromeOptions;
import org.openqa.selenium.remote.RemoteWebDriver;

import java.io.File;
import java.io.IOException;
import java.net.URI;
import java.net.URL;
import java.nio.file.Files;

public class Test {

    public static By byText(String text) {
        return By.xpath ("//*[contains(text(),'"+text+"')]");
    }


    public static void saveFile(WebDriver driver, String where) {
        File scrFile = ((TakesScreenshot)driver).getScreenshotAs(OutputType.FILE);
        File f = new File(where);
        try {
            Files.move(scrFile.toPath(), f.toPath());
            scrFile.delete();
        } catch (IOException e1) {
            throw new RuntimeException(e1);
        }
        if (!f.exists())
            throw new RuntimeException("Target file "+ f.getAbsolutePath() +" was reported created by does not exists!");
//        System.out.println("Wrote file: " + f.getAbsolutePath());
    }

    public static void main(String[] args) {
        WebDriver driver = null;
        try {
            String selenium = "http://selenium:4444/wd/hub";
            String test = "http://maincontainer:8080/page1";
            if (args.length > 0)
                selenium = args[0];
            if (args.length > 1)
                test = args[1];
            // Connect to the Selenium server running in Docker
            URL url = new URI(selenium).toURL();
//            DesiredCapabilities capabilities = DesiredCapabilities.();
            driver = new RemoteWebDriver(url, new ChromeOptions());

            driver.get(test);

            // Find the search box and perform a search
            WebElement bt = driver.findElement(byText("Button2"));
            if (!driver.findElements(byText("Button 3")).isEmpty()) {
                throw new Exception("Button 3 already exists");
            }
            bt.click();
            WebElement bt1 = driver.findElement(byText("Button 3"));
            bt1.click();
            System.out.println("Passed!");
            saveFile(driver, "screenshots/test1_completed.png");    
            driver.quit();
            driver = null;
        } catch (Exception e) {
            e.printStackTrace();
            if (driver != null) 
                saveFile(driver, "screenshots/error.png");    
            System.exit(1); // error state.
        }
    }
}
