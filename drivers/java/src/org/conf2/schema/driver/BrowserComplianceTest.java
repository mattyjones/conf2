package org.conf2.schema.driver;

import org.conf2.schema.DataType;
import org.conf2.schema.HasDataType;
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

    BrowseValue enumValue(int eid) {
        BrowseValue v = new BrowseValue();
        DataType type = ((HasDataType) selection.position).getDataType();
        v.setEnum(type, eid);
        return v;
    }

    BrowseValue enumValues(int... eids) {
        BrowseValue v = new BrowseValue();
        DataType type = ((HasDataType) selection.position).getDataType();
        v.setEnumList(type, eids);
        return v;
    }

    @Test
    public void readTest() {
        readTest("read one a", () -> {
            return enumValue(0);
        });
        readTest("read one b", () -> {
            return enumValue(1);
        });
        readTest("read two a", () -> {
            return enumValues(0);
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
