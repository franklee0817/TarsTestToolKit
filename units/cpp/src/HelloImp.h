#ifndef _HelloImp_H_
#define _HelloImp_H_

#include "servant/Application.h"
#include "Hello.h"

/**
 *
 *
 */
class HelloImp : public Demo::Hello
{
public:
    /**
     *
     */
    virtual ~HelloImp() {}

    /**
     *
     */
    virtual void initialize();

    /**
     *
     */
    virtual void destroy();

    /**
     *
     */
    virtual string ping(const string &req, tars::TarsCurrentPtr current) { return "pong"; };
};
/////////////////////////////////////////////////////
#endif
