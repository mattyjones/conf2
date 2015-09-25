package org.conf2.schema.driver;

import org.conf2.schema.*;
import org.conf2.schema.browse.*;
import org.conf2.schema.comm.DataReader;
import org.conf2.schema.comm.DataWriter;

import java.io.ByteArrayOutputStream;
import java.io.IOException;
import java.nio.ByteBuffer;

/**
 *
 */
public class BrowserAdaptor {
    private final static BrowseValue[] NO_KEYS = new BrowseValue[0];

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

    public static boolean iterate(Selection s, ByteBuffer encodedKeyValues, boolean first) {
        BrowseValue[] key;
        if (encodedKeyValues != null) {
            DataReader r = new DataReader(encodedKeyValues);
            DataType[] keyTypes = ((MetaList)(s.meta)).getKeyDataTypes();
            key = r.readValues(keyTypes);
        } else {
            key = NO_KEYS;
        }
        return s.Iterate.Iterate(key, first);
    }

    public static ByteBuffer read(Selection s, String ident) {
        s.position = updatePosition(s.meta, ident);
        BrowseValue val = s.Read.Read();
        ByteBuffer data = null;
        if (val != null) {
            DataWriter w = new DataWriter();
            w.writeValue(val);
            data = w.getByteBuffer();
        }
        return data;
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

    public static void edit(Selection s, String ident, int opCode, ByteBuffer encodedValue) {
        s.position = updatePosition(s.meta, ident);
        EditOperation op = EditOperation.values()[opCode];
        BrowseValue val = null;
        if (encodedValue != null) {
            DataReader r = new DataReader(encodedValue);
            val = r.readValue(((HasDataType) s.position).getDataType());
        }
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
