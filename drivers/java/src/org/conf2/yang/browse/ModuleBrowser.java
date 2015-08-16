package org.conf2.yang.browse;

import org.conf2.yang.Meta;
import org.conf2.yang.MetaError;
import org.conf2.yang.Module;
import org.conf2.yang.YangModule;

import java.lang.reflect.Method;

/**
 *
 */
public class ModuleBrowser implements Browser {
    Module yang = YangModule.YANG;
    public Module module;

    public ModuleBrowser(Module module) {
        this.yang = yang;
        this.module = module;
    }

    @Override
    public Selection getRootSelector() {
        Selection s = new Selection();
        s.meta = yang;
        s.Enter = () -> enterModule(module);
        s.Edit = new BrowseEdit() {
            @Override
            public void Edit(EditOperation op, BrowseValue val) {
                switch (op) {
                    case CREATE_CHILD:
                        module = new Module("unknown");
                        break;
                    case UPDATE_VALUE:
                        setterField(s.position, module, val);
                        break;
                }
            }
        };
        return s;
    }

    Selection enterModule(Module module) {
        Selection s = new Selection();
        s.found = (module != null);
        s.Read = (BrowseValue v) -> getterField(s.position, module, v);
        return s;
    }

    static String getterMethodNameFromMeta(String ident) {
        return "get" + Character.toUpperCase(ident.charAt(0)) + ident.substring(1);
    }

    static void setterField(Meta m, Object o, BrowseValue v) {

    }

    static void getterField(Meta m, Object o, BrowseValue v) {
        String methodName = getterMethodNameFromMeta(m.getIdent());
        try {
            Method method = o.getClass().getMethod(methodName);
            Object result = method.invoke(o);
            v.str = result.toString();
        } catch (ReflectiveOperationException e) {
            throw new MetaError("Method not found", e);
        }
    }
}
