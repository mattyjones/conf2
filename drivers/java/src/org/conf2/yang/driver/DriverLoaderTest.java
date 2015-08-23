package org.conf2.yang.driver;

import org.conf2.yang.SimpleDataSource;
import org.conf2.yang.Module;
import org.junit.Test;

import static org.junit.Assert.*;

public class DriverLoaderTest {

    @Test
    public void testLoadModule() {
        Driver d = new Driver();
	DriverLoader loader = new DriverLoader(d);
	SimpleDataSource ds = new SimpleDataSource(DriverLoaderTest.class);
	Module simple = loader.loadModule(ds, "testdata/simple.yang");
	assertNotNull(simple);
	assertEquals("turing-machine", simple.getIdent());
    }
}
