package org.conf2.schema.browse;

import org.conf2.schema.*;
import org.conf2.schema.driver.Driver;
import org.conf2.schema.driver.DriverHandle;

import java.util.Collections;
import java.util.Iterator;

/**
 *
 */
public class Walk {
//    private Driver driver;
//    public void insert(Selection from, Selection to) {
//        DriverHandle dfrom = getSelectionHandle(from);
//        DriverHandle dto = getSelectionHandle(from);
//
//
//    }
//
//
//    DriverHandle getSelectionHandle(Browser b, String path, Selection s) {
//
//
//    }

    public static Selection walk(Selection selection, BrowsePath path) {
        PathWalkController controller = new PathWalkController(path);
        walk(selection, controller);
        return controller.target;
    }

    public static void walk(Selection selection, WalkController controller) {
        Walk.walk(selection, controller, 0);
    }

    static void walk(Selection s, WalkController c, int level) {
        if (s.meta instanceof MetaList && ! s.insideList) {
            MetaList l = (MetaList) s.meta;
            boolean hasMore = c.selectionIterator(s, level, true);
            while (hasMore) {
                s.insideList = true;
                walk(s, c, level);
                hasMore = c.selectionIterator(s, level, false);
            }
        } else {
            Selection child;
            Iterator<Meta> i = c.containerIterator(s, level);
            while (i.hasNext()) {
                s.position = i.next();

                // TODO: resolve choice

                if (MetaUtil.isLeaf(s.position)) {
                    BrowseValue v = new BrowseValue();
                    // just walking, not interested in value
                    s.Read.Read();
                } else {
                    child = s.Enter.Enter();
                    if (!s.found) {
                        continue;
                    }
                    child.meta = (MetaCollection) s.position;
                    walk(child, c, level + 1);

                    if (s.Exit != null) {
                        s.Exit.Exit();
                    }
                }
            }
        }
    }
}


