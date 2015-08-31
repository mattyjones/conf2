package org.conf2.yang.driver;

import org.conf2.yang.MetaUtil;
import org.conf2.yang.SimpleStreamSource;
import org.conf2.yang.Module;
import org.junit.Test;

import java.io.IOException;

import static org.junit.Assert.*;

public class DriverLoaderTest {

    @Test
    public void testLoadModule() throws IOException {
        Driver d = Driver.yangDriver();
		SimpleStreamSource ds = new SimpleStreamSource(DriverLoaderTest.class);
		Module simple = d.loadModule(ds, "simple.yang");
		assertNotNull(simple);
		assertEquals("turing-machine", simple.getIdent());
		assertEquals(4, MetaUtil.collectionLength(simple.getTypedefs()));
		assertEquals(1, MetaUtil.collectionLength(simple.getGroupings()));
		assertEquals(1, MetaUtil.collectionLength(simple));
		assertEquals(2, MetaUtil.collectionLength(simple.getRpcs()));
		assertEquals(1, MetaUtil.collectionLength(simple.getNotifications()));
	}
}
