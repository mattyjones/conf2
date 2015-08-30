package org.conf2.yang.driver;

import org.conf2.yang.*;
import org.conf2.yang.browse.BrowseValue;
import org.conf2.yang.browse.Browser;
import org.conf2.yang.browse.EditOperation;
import org.conf2.yang.browse.Selection;

/**
 *
 */
public class BrowserAdaptor {
    private final static String[] NO_KEYS = new String[0];

    public static Selection getRootSelector(Browser browser) {
        return browser.getRootSelector();
    }

    public static Selection enter(Selection s, String ident) {
        s.position = updatePosition(s.meta, ident);
        Selection child = s.Enter.Enter();
        if (child != null) {
            child.meta = (MetaCollection) s.position;
        }
        return child;
    }

    public static boolean iterate(Selection s, String encodedKeys, boolean first) {
        String[] keys;
        if (encodedKeys != null) {
            keys = encodedKeys.split(" ");
        } else {
            keys = NO_KEYS;
        }
        return s.Iterate.Iterate(keys, first);
    }

    public static BrowseValue read(Selection s, String ident) {
        s.position = updatePosition(s.meta, ident);
        BrowseValue val = new BrowseValue();
        s.Read.Read(val);
        return val;
    }

    static Meta updatePosition(MetaCollection meta, String ident) {
        Meta position = null;
        if (ident != null) {
            int choiceSep = ident.indexOf('/');
            if (choiceSep >= 0) {
                String choiceIdent = ident.substring(0, choiceSep);
                Choice c = (Choice) MetaUtil.findByIdent(meta, choiceIdent);
                if (c != null) {
                    int caseSep = ident.indexOf('/', choiceSep + 1);
                    String caseIdent = ident.substring(choiceSep + 1, caseSep);
                    ChoiceCase ccase = c.getCase(caseIdent);
                    if (ccase != null) {
                        position = new MetaCollectionIterator(ccase).next();
                    }
                }
            } else {
                position = MetaUtil.findByIdent(meta, ident);
            }
            if (position == null) {
                throw new DriverError("Could not find meta " + ident + " in " + meta.getIdent());
            }
        }
        return position;
    }

    public static void edit(Selection s, String ident, int opCode, BrowseValue val) {
        s.position = updatePosition(s.meta, ident);
        EditOperation op = EditOperation.values()[opCode];
        s.Edit.Edit(op, val);
    }

    public static void exit(Selection s, String ident) {
        s.position = updatePosition(s.meta, ident);
        if (s.Exit != null) {
            s.Exit.Exit();
        }
    }

    public static String choose(Selection s, String choiceIdent) {
        Choice choice = (Choice) MetaUtil.findByIdent(s.meta, choiceIdent);
        s.position = s.Choose.Choose(choice);
        return s.position.getIdent();
    }
}
