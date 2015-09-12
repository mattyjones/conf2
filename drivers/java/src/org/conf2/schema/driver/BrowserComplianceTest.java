package org.conf2.schema.driver;

import org.conf2.schema.ValueType;
import org.conf2.schema.browse.*;
import org.junit.AfterClass;
import org.junit.Before;
import org.junit.BeforeClass;
import org.junit.Test;
import static org.junit.Assert.*;

import java.util.HashMap;
import java.util.Map;

/**
 *
 */
public class BrowserComplianceTest {
    static Driver driver;
    static DriverTestHarness harness;
    boolean passed;
    Selection selection;

    @BeforeClass
    public static void setUpAll() {
        driver = new Driver();
        harness = new DriverTestHarness(driver);
    }

    @AfterClass
    public static void teardownAll() {
        harness.release();
    }

    @Before
    public void setUp() {
        passed = true;
        selection = new Selection();
    }

    @Test
    public void readTest() {
        readTest("read one a", (BrowseValue v) -> {
            v.setEnum(selection.position, 0);
        });
        readTest("read one b", (BrowseValue v) -> {
            v.setEnum(selection.position, 1);
        });
        readTest("read two a", (BrowseValue v) -> {
            v.addEnum(selection.position, 0);
        });
        System.out.print(harness.getReport());
        assertTrue(passed);
    }

    @Test
    public void editTest() {
        writeTest("write one a", (EditOperation op, BrowseValue val) -> {
            System.out.printf("Here");
        });
    }

    void readTest(String testname, BrowseRead read) {
        selection.Read = read;
        passed = harness.runTest(testname, selection) && passed;
    }

    void writeTest(String testname, BrowseEdit edit) {
        selection.Edit = edit;
        selection.found = true;
        passed = harness.runTest(testname, selection) && passed;
    }
}
