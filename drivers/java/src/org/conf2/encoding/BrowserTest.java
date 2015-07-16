package org.conf2.encoding;

import org.junit.Before;
import org.junit.Test;

import java.io.ByteArrayOutputStream;
import java.io.IOException;
import java.net.MalformedURLException;
import java.util.Map;

import static org.junit.Assert.*;

/**
 *
 */
public class BrowserTest {
    private Browser db;
    private Schema root;
    private Schema animalNode;
    private Schema birdNode;
    private DataRoot dataRoot = new DataRoot();
    private SelectorHandler selectorHandler;
    private AssemblerFactory factory;
    private Map<String,Object> birds;
    private ByteArrayOutputStream actual;

    static class DataRoot implements Assembler {
        Animal animal = new Animal();

        @Override
        public void write(Renderer wtr) throws IOException {
            wtr.enterContainer("animal");
            animal.write(wtr);
            wtr.leaveContainer();
        }
    }

    static class Bird implements Assembler {
        Object canary = "tweet";
        Object pigeon = "coo";

        public void write(Renderer wtr) throws IOException {
            wtr.writeValue("canary",canary);
            wtr.writeValue("pigeon", pigeon);
        }
    }

    static class Animal implements Assembler {
        Bird bird = new Bird();

        public void write(Renderer wtr) throws IOException {
            wtr.enterContainer("bird");
            bird.write(wtr);
            wtr.leaveContainer();
        }
    }

    @Before
    public void newTestDb() throws MalformedURLException {
        root = new Schema("");
        animalNode = root.addChild(new Schema("animal"));
        birdNode = animalNode.addChild(new Schema("bird"));
        db = new Browser(root, dataRoot);

        selectorHandler = new SelectorHandler() {
            @Override
            public Object nextLevel(Schema node, String element, Object data) {
                if (node.getIdent().equals("animal")) {
                    return dataRoot.animal;
                } else if (node.getIdent().equals("bird")) {
                    return dataRoot.animal.bird;
                }
                return null;
            }
        };
        factory = new AssemblerFactory() {
            @Override
            public Assembler getAssembler(Object obj) {
                return (Assembler)obj;
            }
        };
        actual = new ByteArrayOutputStream();
    }

    @Test
    public void testSearch() throws IOException {
        Object results = db.search(new Selector("animal/bird"), selectorHandler);
        assertSame(dataRoot.animal.bird, results);
        JsonRenderer wtr = new JsonRenderer(actual);
        wtr.start();
        factory.getAssembler(results).write(wtr);
        wtr.end();
        assertEquals("{\"canary\":\"tweet\",\"pigeon\":\"coo\"}", new String(actual.toByteArray()));

        results = db.search(new Selector("animal"), selectorHandler);
        assertSame(dataRoot.animal, results);
        actual.reset();
        wtr = new JsonRenderer(actual);
        wtr.start();
        factory.getAssembler(results).write(wtr);
        wtr.end();
        assertEquals("{\"bird\":{\"canary\":\"tweet\",\"pigeon\":\"coo\"}}", new String(actual.toByteArray()));
    }

    @Test
    public void testDepth() throws IOException {
    }
}
