package org.conf2.yang.driver;

import org.conf2.yang.*;
import org.junit.Test;
import static org.junit.Assert.*;

/**
 *
 */
public class BrowserAdaptorTest {
    @Test
    public void testUpdatePosition() {
        Container container = new Container("x");
        Choice choice = new Choice("c");
        container.addMeta(choice);
        ChoiceCase case1 = new ChoiceCase("y");
        Leaf l1 = new Leaf("l1");
        case1.addMeta(l1);
        choice.addMeta(case1);
        ChoiceCase case2 = new ChoiceCase("z");
        Leaf l2 = new Leaf("l2");
        case1.addMeta(l2);
        choice.addMeta(case2);

        Meta actual = BrowserAdaptor.updatePosition(container, "c/y/l1");
        assertSame(actual, l1);
    }
}
