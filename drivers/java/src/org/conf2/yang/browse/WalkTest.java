package org.conf2.yang.browse;

import org.conf2.yang.YangModule;
import org.junit.Test;

/**
 *
 */
public class WalkTest {

    @Test
    public void testExhaustiveWalk() {
        ModuleBrowser browser = new ModuleBrowser(YangModule.YANG);
        Walk.walk(browser.getRootSelector(), new ExhaustiveWalk());
    }
}
