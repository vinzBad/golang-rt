use strict;
use warnings;

our @Users = ({
    Name         => 'john.doe',
    Password     => 'password',
    EmailAddress => 'john.doe@localhost',
    Privileged   => 1,
    Disabled     => 0,
    },
    {
    Name         => 'foo.bar',
    Password     => 'password',
    EmailAddress => 'foo.bar@localhost',
    Privileged   => 1,
    Disabled     => 0,
    },
);


our @ACL = ();

push @ACL, map {
    {
        Right       => $_,
        GroupDomain => 'SystemInternal',
        Queue       => 'General',
        GroupType   => 'Privileged',
    }
} qw(ShowTicket ReplyToTicket CreateTicket Watch SeeQueue);

1;

