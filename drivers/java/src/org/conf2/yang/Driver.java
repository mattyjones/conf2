package org.conf2.yang;

public class Driver {
    private static boolean init;
    static {
        System.loadLibrary("yangj");
    }

    public Driver() {
        if (!init) {
            initializeDriver();
            init = true;
        }
    }

    private native void initializeDriver();
}
