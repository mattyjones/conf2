package org.conf2.yang.driver;

import org.conf2.yang.DataSource;
import org.junit.Test;

import java.io.ByteArrayInputStream;
import java.io.IOException;
import java.io.InputStream;

import static org.junit.Assert.assertEquals;

public class DriverTest {

    @Test
    public void testLoadResource() throws IOException {
        final String expectedResourceContent = "resourceContentHere";
        final String[] actualResourceId = new String[1];
        final DataSource r = new DataSource() {

            @Override
            public InputStream getResource(String resourceId) throws IOException {
                actualResourceId[0] = resourceId;
                return new ByteArrayInputStream(expectedResourceContent.getBytes());
            }
        };
        String expectedResourceId = "resourceIdHere";
        Driver d = new Driver();
        String actualResource = d.echoTest(r, expectedResourceId);

        assertEquals(expectedResourceId, actualResourceId[0]);
        assertEquals(expectedResourceContent, actualResource);
    }
}
