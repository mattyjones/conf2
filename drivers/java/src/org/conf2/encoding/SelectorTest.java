package org.conf2.encoding;

import org.junit.Test;

import java.net.MalformedURLException;

import static org.junit.Assert.assertArrayEquals;
import static org.junit.Assert.assertEquals;

/**
 *
 */
public class SelectorTest {

    @Test
    public void testElements() throws MalformedURLException {
        Selector p = new Selector("blah/flaw=glop/glaw?bleep/bloop");
        assertEquals(3, p.getElements().length);
        String[] expected = {
                "blah",
                "flaw=glop",
                "glaw"
        };
        assertArrayEquals(expected, p.getElements());
    }

}
