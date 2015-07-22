package org.conf2.yang.comm;

import org.conf2.yang.*;

import java.lang.reflect.Constructor;
import java.lang.reflect.InvocationTargetException;
import java.lang.reflect.Method;
import java.util.HashMap;
import java.util.Map;
import java.util.Stack;

interface Reciever {
    void startTransaction();
    void newContainer(Container m);
    void newList(MetaList m);
    void newListItem(MetaList m);
    void putIntLeaf(Leaf m, int v);
    void putStringLeaf(Leaf m, String v);
    void putStringLeafList(LeafList m, String v);
    void putIntLeafList(LeafList m, int[] v);
    void exitContainer(Container m);
    void exitList(MetaList m);
    void exitListItem(MetaList m);
    void endTransaction();
}

/**
 *
 */
public class MetaReciever implements Reciever {
    private Module module;
    private Stack<Meta> stack;
    private Constructor metaClass;
    static Map<String,Constructor> metaFactory = new HashMap<String,Constructor>();
    static Map<String,Method> reflectionCache = new HashMap<String, Method>();
    static {
        Class identParam[] = new Class[] {String.class};
        try {
            metaFactory.put("module", Module.class.getConstructor(identParam));
            metaFactory.put("container", Container.class.getConstructor(identParam));
            metaFactory.put("list", MetaList.class.getConstructor(identParam));
            metaFactory.put("leaf", Leaf.class.getConstructor(identParam));
            metaFactory.put("leaf-list", LeafList.class.getConstructor(identParam));
            metaFactory.put("grouping", Grouping.class.getConstructor(identParam));
            metaFactory.put("choice", Choice.class.getConstructor(identParam));
            metaFactory.put("case", ChoiceCase.class.getConstructor(identParam));
        } catch (NoSuchMethodException e) {
            throw new RuntimeException(e);
        }
    }
    static Class[] SET_STRING = new Class[] {String.class};
    static Class[] SET_INT = new Class[] {Integer.class};

    @Override
    public void startTransaction() {
        stack = new Stack<Meta>();
    }

    @Override
    public void newContainer(Container m) {
        if (m.getIdent() == "module") {
            metaClass = metaFactory.get(m.getIdent());
        }
    }

    @Override
    public void newList(MetaList m) {
        stack.push(m);
    }

    @Override
    public void newListItem(MetaList m) {
        stack.push(m);
    }

    @Override
    public void putIntLeaf(Leaf m, int v) {
        Meta data = stack.peek();
        try {
            Method setter = getMethod(data.getClass(), m.getIdent(), SET_INT);
            setter.invoke(data, v);
        } catch (ReflectiveOperationException e) {
            throw new MetaError("Could not set value " + m.getIdent(), e);
        }
    }

    Method getMethod(Class dataClass, String setter, Class[] params) throws NoSuchMethodException {
        String key = dataClass.getCanonicalName() + setter;
        Method m = reflectionCache.get(key);
        if (m == null) {
            m = dataClass.getMethod(setter, params);
            reflectionCache.put(key, m);
        }
        return m;
    }

    @Override
    public void putStringLeaf(Leaf m, String v) {
        if (m.getIdent().equals("ident")) {
            if (metaClass == null) {
                throw new MetaError("No meta class defined");
            }
            try {
                Meta data = (Meta) metaClass.newInstance(v);
                stack.push(data);
            } catch (ReflectiveOperationException e) {
                throw new MetaError("Could not construct meta object " + v, e);
            }
        } else {
            Meta data = stack.peek();
            try {
                Method setter = getMethod(data.getClass(), m.getIdent(), SET_INT);
                setter.invoke(data, v);
            } catch (ReflectiveOperationException e) {
                throw new MetaError("Could not set value " + m.getIdent(), e);
            }
        }
    }

    @Override
    public void putStringLeafList(LeafList m, String v) {
        // TODO
    }

    @Override
    public void putIntLeafList(LeafList m, int[] v) {
        // TODO
    }

    @Override
    public void exitContainer(Container m) {
        stack.pop();
    }

    @Override
    public void exitList(MetaList m) {
        stack.pop();
    }

    @Override
    public void exitListItem(MetaList m) {
        stack.pop();
    }

    @Override
    public void endTransaction() {
    }
}
