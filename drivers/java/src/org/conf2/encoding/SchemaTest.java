package org.conf2.encoding;

import org.junit.Test;

import java.net.MalformedURLException;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertSame;

/**
 *
 */
public class SchemaTest {
    private Schema root;
    private Schema animalNode;
    private Schema birdNode;

    public void newTestDb() throws MalformedURLException {
        root = new Schema("");
        animalNode = root.addChild(new Schema("animal"));
        birdNode = animalNode.addChild(new Schema("bird"));
    }

    @Test
    public void testConstruction() {
        Schema parent = new Schema("root");
        Schema a1 = parent.addChild(new Schema("a1"));
        parent.addChild(new Schema("a2"));
        a1.addChild(new Schema("b1"));
        assertEquals(2, parent.size());
    }

    @Test
    public void testGetChildByIdent() {
        assertSame(animalNode, root.getChildByIdent("animal"));

    }
}
